package sched

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/btree"
)

func (s *cronScheduler) signal() {
	select {
	case s.signalCh <- struct{}{}:
	default:
	}
}

func (s *cronScheduler) Start(ctx context.Context) {
	go s.run(ctx)
}

var MaxTime = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

func (s *cronScheduler) run(ctx context.Context) {
	duration := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.signalCh:
		case <-time.After(duration):
		}

		nextTick := s.onTick()
		duration = time.Until(nextTick)
		if nextTick.Equal(MaxTime) {
			duration = time.Hour
		}

		log.WithField("nextTick", time.Now().Add(duration).Truncate(time.Second)).
			Info("scheduling next trigger")
	}
}

type CronScheduler interface {
	Start(ctx context.Context)
	Schedule(id int64, at time.Time)
	Remove(id int64) error
}

type cronScheduler struct {
	mtx sync.RWMutex

	index      *btree.BTree
	signalCh   chan struct{}
	onCronTick func(id int64) time.Time
}

func NewCronScheduler(onCronTick func(cronID int64) time.Time) CronScheduler {
	return &cronScheduler{
		index:      btree.New(64),
		signalCh:   make(chan struct{}, 1),
		onCronTick: onCronTick,
	}
}

func (s *cronScheduler) Schedule(id int64, at time.Time) {
	s.mtx.Lock()

	s.index.ReplaceOrInsert(&item{
		id:         id,
		nextTickAt: at.Unix(),
	})

	s.mtx.Unlock()

	s.signal()
}

type item struct {
	id         int64
	nextTickAt int64
}

func (i *item) Less(than btree.Item) bool {
	other := than.(*item)
	if i.nextTickAt == other.nextTickAt {
		return i.id < other.id
	}
	return i.nextTickAt < other.nextTickAt
}

func (s *cronScheduler) onTick() time.Time {
	now := time.Now().Unix()

	s.mtx.Lock()
	defer s.mtx.Unlock()

	for {
		it, _ := s.index.Min().(*item)
		if it == nil {
			return MaxTime
		}

		if it.nextTickAt > now {
			return time.Unix(it.nextTickAt, 0)
		}

		s.index.DeleteMin()

		nextTick := s.onCronTick(it.id)

		if nextTick.Unix() > now {
			s.index.ReplaceOrInsert(&item{
				id:         it.id,
				nextTickAt: nextTick.Unix(),
			})
		}
	}
}

func (s *cronScheduler) Remove(id int64) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.index.Delete(&item{id: id})

	return nil
}
