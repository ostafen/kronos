package sched

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/suite"
)

type ScheduleServiceSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, &ScheduleServiceSuite{})
}

func (s *ScheduleServiceSuite) SetupSuite() {
	rand.Seed(time.Now().Unix())

}

func (s *ScheduleServiceSuite) TestOnTickIsCalledProperly() {
	now := time.Now()

	schedules := make(map[string]bool)

	calls := 0
	manager := NewScheduleManager(func(id string) time.Time {
		s.True(schedules[id])

		calls++
		return time.Time{}
	})

	n := 1000
	for i := 0; i < n; i++ {
		id := uuid.NewString()
		schedules[id] = true
		manager.Schedule(id, now.Add(-time.Hour*time.Duration(i)))
		manager.Schedule(uuid.NewString(), now.Add(time.Hour*time.Duration(i+1)))
	}

	s.Len(schedules, n)
	s.Equal(manager.(*scheduleManager).index.Len(), 2*n)

	manager.(*scheduleManager).onTick()
	s.Equal(calls, n)
}

func (s *ScheduleServiceSuite) TestOnTickReschedule() {
	n := 1000

	now := time.Now().Truncate(time.Second)

	rescheduled := 0
	manager := NewScheduleManager(func(id string) time.Time {
		if rand.Int()%2 == 0 {
			rescheduled++
			return now.Add(time.Second * time.Duration(rescheduled))
		}
		return time.Time{}
	})

	for i := 0; i < n; i++ {
		manager.Schedule(uuid.NewString(), now)
	}

	m := manager.(*scheduleManager)
	m.onTick()

	s.Equal(m.index.Len(), rescheduled)

	var expectedNextTick time.Time
	if rescheduled > 0 {
		expectedNextTick = now.Add(time.Second)
	}
	s.Equal(expectedNextTick, m.nextTick())
}
