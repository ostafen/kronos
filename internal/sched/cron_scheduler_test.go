package sched

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CronSchedulerServiceSuite struct {
	suite.Suite
}

func TestCronScheduleServiceSuite(t *testing.T) {
	suite.Run(t, &CronSchedulerServiceSuite{})
}

func (s *CronSchedulerServiceSuite) SetupSuite() {
	rand.Seed(time.Now().Unix())
}

func (s *CronSchedulerServiceSuite) TestOnTickIsCalledProperly() {
	now := time.Now()

	schedules := make(map[int64]bool)

	calls := 0
	scheduler := NewCronScheduler(func(id int64) time.Time {
		s.True(schedules[id])

		calls++
		return time.Time{}
	})

	n := 1000
	for i := 0; i < n; i++ {
		id := rand.Int63()
		schedules[id] = true
		scheduler.Schedule(id, now.Add(-time.Hour*time.Duration(i)))
		scheduler.Schedule(rand.Int63(), now.Add(time.Hour*time.Duration(i+1)))
	}

	s.Len(schedules, n)
	s.Equal(scheduler.(*cronScheduler).index.Len(), 2*n)

	scheduler.(*cronScheduler).onTick()
	s.Equal(calls, n)
}

func (s *CronSchedulerServiceSuite) TestOnTickReschedule() {
	n := 1000

	now := time.Now().Truncate(time.Second)

	rescheduled := 0
	scheduler := NewCronScheduler(func(id int64) time.Time {
		if rand.Int()%2 == 0 {
			rescheduled++
			return now.Add(time.Second * time.Duration(rescheduled))
		}
		return time.Time{}
	})

	for i := 0; i < n; i++ {
		scheduler.Schedule(rand.Int63(), now)
	}

	m := scheduler.(*cronScheduler)
	nextTick := m.onTick()

	s.Equal(m.index.Len(), rescheduled)

	var expectedNextTick time.Time
	if rescheduled > 0 {
		expectedNextTick = now.Add(time.Second)
	}
	s.Equal(expectedNextTick, nextTick)
}
