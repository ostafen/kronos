package notification

import (
	"context"
	"time"

	"github.com/ostafen/kronos/internal/service"
	log "github.com/sirupsen/logrus"
)

type scheduleTrigger struct {
	schedService service.ScheduleService
	signal       chan struct{}
}

func NewScheduleTrigger(schedService service.ScheduleService) *scheduleTrigger {
	return &scheduleTrigger{
		schedService: schedService,
		signal:       make(chan struct{}, 1),
	}
}

func (t *scheduleTrigger) WakeUp() {
	t.signal <- struct{}{}
}

func (t *scheduleTrigger) Start(ctx context.Context) {
	go t.run(ctx)
}

const (
	MaxSleepDuration = time.Hour
)

func (t *scheduleTrigger) run(ctx context.Context) {
	duration := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.signal:
		case <-time.After(duration):
		}

		nextScheduleTime, err := t.schedService.NotifySchedules()
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now()
		if nextScheduleTime == nil {
			duration = MaxSleepDuration
		} else if nextScheduleTime.After(now) {
			duration = time.Until(*nextScheduleTime)
		} else {
			duration = time.Duration(0)
		}

		log.WithField("nextTick", time.Now().Add(duration).Truncate(time.Second)).
			Info("scheduling next trigger")
	}
}
