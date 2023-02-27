package sched

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/btree"
)

func (s *scheduleManager) WakeUp() {
	s.signal <- struct{}{}
}

func (s *scheduleManager) Start(ctx context.Context) {
	go s.run(ctx)
}

const (
	MaxTime = 1<<63 - 1
)

func (s *scheduleManager) run(ctx context.Context) {
	duration := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.signal:
		case <-time.After(duration):
		}

		nextTick := s.onTick()
		duration = time.Until(nextTick)
		if nextTick.Unix() == MaxTime {
			duration = time.Hour
		}

		log.Println(duration)
		log.WithField("nextTick", time.Now().Add(duration).Truncate(time.Second)).
			Info("scheduling next trigger")
	}
}

type ScheduleManager interface {
	Start(ctx context.Context)
	WakeUp()
	Schedule(id string, at time.Time)
	Remove(id string) error
}

type scheduleManager struct {
	mtx sync.RWMutex

	deletedIds map[string]struct{}
	index      *btree.BTree
	signal     chan struct{}
	iterFunc   func(id string) time.Time
}

const (
	btreeDegree = 2
)

func NewScheduleManager(iterFunc func(id string) time.Time) ScheduleManager {
	s := &scheduleManager{
		index:      btree.New(btreeDegree),
		signal:     make(chan struct{}, 1),
		iterFunc:   iterFunc,
		deletedIds: make(map[string]struct{}),
	}
	return s
}

func (s *scheduleManager) Schedule(id string, at time.Time) {
	s.mtx.Lock()

	delete(s.deletedIds, id)

	s.index.ReplaceOrInsert(&item{
		id:                id,
		nextTickTimestamp: at.Unix(),
	})

	s.mtx.Unlock()
}

type item struct {
	id                string
	nextTickTimestamp int64
}

func (i *item) Less(than btree.Item) bool {
	other := than.(*item)
	if i.nextTickTimestamp == other.nextTickTimestamp {
		return i.id < other.id
	}
	return i.nextTickTimestamp < other.nextTickTimestamp
}

func (s *scheduleManager) onTick() time.Time {
	rescheduleItems := make([]*item, 0)
	removeItems := make([]*item, 0)

	now := time.Now().Unix()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.index.Ascend(func(i btree.Item) bool {
		it, _ := i.(*item)
		if it == nil {
			return false
		}

		if _, deleted := s.deletedIds[it.id]; deleted {
			removeItems = append(removeItems, it)
			return true
		}

		include := it.nextTickTimestamp <= now
		if include {
			nextTickTime := s.iterFunc(it.id)
			removeItems = append(removeItems, it)

			nextTick := nextTickTime.Unix()
			if nextTick > now {
				rescheduleItems = append(rescheduleItems, &item{
					id:                it.id,
					nextTickTimestamp: nextTickTime.Unix(),
				})
			}
		}
		return include
	})

	for _, it := range removeItems {
		if s.index.Delete(it) == nil {
			return time.Time{}
		}
	}

	for _, it := range rescheduleItems {
		s.index.ReplaceOrInsert(it)
	}

	return s.nextTick()
}

func (s *scheduleManager) nextTick() time.Time {
	it, _ := s.index.Min().(*item)
	if it == nil {
		return time.Unix(1<<63-1, 0)
	}
	return time.Unix(it.nextTickTimestamp, 0)
}

func (s *scheduleManager) Remove(id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.deletedIds[id] = struct{}{}
	return nil
}
