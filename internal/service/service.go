package service

import (
	"context"
	"errors"
	"time"

	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/sched"
	"github.com/ostafen/kronos/internal/store"

	log "github.com/sirupsen/logrus"
)

type ScheduleService interface {
	RegisterSchedule(sched *model.ScheduleRegisterInput) (*model.CronSchedule, error)
	GetSchedule(id int64) (*model.CronSchedule, error)
	DeleteSchedule(id int64) error
	IterSchedules(onSchedule func(*model.CronSchedule) error) error
	GetHistory() ([]*model.CronStatus, error)
	GetCronHistory(cronID int64) ([]*model.CronStatus, error)

	PauseSchedule(id int64) (*model.CronSchedule, error)
	ResumeSchedule(id int64) (*model.CronSchedule, error)
	TriggerSchedule(id int64) (*model.CronSchedule, error)

	Scheduler() sched.CronScheduler
	Stop()
}

func NewScheduleService(
	store store.Store,
	notificationSvc NotificationService,
) ScheduleService {
	svc := &schedService{
		cronRepo:        store.CronScheduleRepository(),
		statusRepo:      store.HistoryRepository(),
		notificationSvc: notificationSvc,
	}
	svc.scheduler = sched.NewCronScheduler(svc.OnTick)

	err := svc.cronRepo.Iter(func(sched *model.CronSchedule) error {
		if sched.IsActive() {
			log.Infof("scheduling %d at %s", sched.ID, sched.NextTick())

			svc.scheduler.Schedule(sched.ID, sched.NextTick())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	svc.cancel = cancel

	svc.scheduler.Start(ctx)
	return svc
}

const MaxSamplesPerCronDefault = 100

type schedService struct {
	notificationSvc NotificationService

	scheduler  sched.CronScheduler
	cronRepo   store.CronScheduleRepository
	statusRepo store.CronHistoryRepository
	cancel     context.CancelFunc
}

func (s *schedService) RegisterSchedule(input *model.ScheduleRegisterInput) (*model.CronSchedule, error) {
	sched, err := input.ToSched()
	if err != nil {
		return nil, err
	}

	id, err := s.cronRepo.Save(sched)
	if err == nil {
		s.scheduler.Schedule(id, sched.NextTick())
	}
	sched.ID = id
	return sched, err
}

const (
	MaxRequestDuration = time.Second * 5
)

func (s *schedService) OnTick(cronID int64) time.Time {
	cron, err := s.cronRepo.Get(cronID)
	if errors.Is(err, store.ErrScheduleNotExist) {
		log.Errorf("no schedule with id %d", cronID)
		return time.Time{}
	}

	if err != nil {
		log.Error(err)
		return time.Time{}
	}

	go func() {
		start := time.Now().Truncate(time.Second)
		status, _ := s.sendWebhookNotification(cron)

		duration := time.Since(start)
		err := s.statusRepo.Insert(&model.CronStatus{
			CronID:     cronID,
			At:         start,
			StatusCode: status,
			Duration:   duration,
		}, MaxSamplesPerCronDefault)

		if err != nil {
			log.Error(err)
		}
	}()

	if cron.Expired() {
		return time.Time{}
	}
	return cron.NextTick()
}

func (s *schedService) sendWebhookNotification(sched *model.CronSchedule) (int, error) {
	log.WithField("scheduleId", sched.ID).
		WithField("url", sched.URL).
		Info("sendingNotification")

	ctx, cancel := context.WithTimeout(context.Background(), MaxRequestDuration)
	defer cancel()
	return s.notificationSvc.Send(ctx, sched.URL, sched)
}

func (s *schedService) GetSchedule(id int64) (*model.CronSchedule, error) {
	sched, err := s.cronRepo.Get(id)
	if err != nil {
		return nil, err
	}
	return sched, nil
}

func (s *schedService) DeleteSchedule(id int64) error {
	return s.cronRepo.Delete(id)
}

func (s *schedService) IterSchedules(onSched func(*model.CronSchedule) error) error {
	return s.cronRepo.Iter(onSched)
}

func (s *schedService) PauseSchedule(id int64) (*model.CronSchedule, error) {
	log.WithField("scheduleId", id).Info("pausing schedule")

	sched, err := s.cronRepo.Get(id)
	if err != nil {
		return nil, err
	}

	sched.Status = model.ScheduleStatusPaused
	if _, err := s.cronRepo.Save(sched); err != nil {
		return nil, err
	}

	s.scheduler.Remove(sched.ID)

	return sched, nil
}

func (s *schedService) TriggerSchedule(id int64) (*model.CronSchedule, error) {
	sched, err := s.cronRepo.Get(id)
	if err != nil {
		return nil, err
	}

	_, err = s.sendWebhookNotification(sched)
	return sched, err
}

func (s *schedService) ResumeSchedule(id int64) (*model.CronSchedule, error) {
	sched, err := s.cronRepo.Get(id)
	if err != nil {
		return nil, err
	}

	sched.Status = model.ScheduleStatusActive
	if _, err := s.cronRepo.Save(sched); err != nil {
		return nil, err
	}

	log.WithField("scheduleId", sched.ID).
		WithField("nextScheduleAt", sched.NextTick()).
		Info("resuming schedule")

	s.scheduler.Schedule(sched.ID, sched.NextTick())

	return sched, nil
}

func (s *schedService) GetCronHistory(cronID int64) ([]*model.CronStatus, error) {
	return s.statusRepo.GetCronHistory(cronID, MaxSamplesPerCronDefault)
}

func (s *schedService) GetHistory() ([]*model.CronStatus, error) {
	return s.statusRepo.GetHistory(MaxSamplesPerCronDefault)
}

func (s *schedService) Scheduler() sched.CronScheduler {
	return s.scheduler
}

func (s *schedService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}
