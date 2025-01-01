package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/store"

	"github.com/stretchr/testify/suite"
)

type ScheduleServiceSuite struct {
	suite.Suite
	schedules           map[string]*model.CronSchedule
	svc                 ScheduleService
	store               store.Store
	webhookHandlerCalls atomic.Int32
}

func TestSuite(t *testing.T) {
	suite.Run(t, &ScheduleServiceSuite{})
}

func (s *ScheduleServiceSuite) SetupTest() {
	s.webhookHandlerCalls.Store(0)

	s.store = &mockStore{}
	s.svc = NewScheduleService(s.store, NewNotificationService())

	s.schedules = make(map[string]*model.CronSchedule)
}

func (s *ScheduleServiceSuite) aSchedule(url string) *model.CronSchedule {
	sched := &model.CronSchedule{
		ID:          rand.Int63(),
		Title:       "test-schedule",
		Status:      model.ScheduleStatusActive,
		CronExpr:    "* * * * * *",
		URL:         url,
		IsRecurring: true,
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(time.Hour),
		CreatedAt:   time.Now(),
		Failures:    0,
	}

	_, err := s.store.CronScheduleRepository().Save(sched)
	s.NoError(err)
	s.svc.Scheduler().Schedule(sched.ID, time.Now())
	return sched
}

func (s *ScheduleServiceSuite) anListeningWebhookHandler(n int, ch chan struct{}) string {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		calls := s.webhookHandlerCalls.Add(1)

		data, err := io.ReadAll(r.Body)
		s.NoError(err)

		var sched model.CronSchedule
		err = json.Unmarshal(data, &sched)
		s.NoError(err)

		_, err = s.store.CronScheduleRepository().Get(sched.ID)
		s.NoError(err)

		if int(calls) == n {
			ch <- struct{}{}
		}
	})
	server := httptest.NewServer(router)
	return fmt.Sprintf("http://%s/webhook", server.Listener.Addr())
}

func (s *ScheduleServiceSuite) TestSchedulesNotification() {
	n := 1000

	ch := make(chan struct{}, 1)

	url := s.anListeningWebhookHandler(n, ch)
	for i := 0; i < n; i++ {
		s.aSchedule(url)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	s.svc.Scheduler().Start(ctx)

	select {
	case <-ctx.Done():
		s.Fail("signal not received")
	case <-ch:
	}
	s.Equal(int(s.webhookHandlerCalls.Load()), n)
}

func (s *ScheduleServiceSuite) TestPauseResume() {
	ch := make(chan struct{}, 1)
	url := s.anListeningWebhookHandler(math.MaxInt, ch)

	sched := s.aSchedule(url)

	time.Sleep(time.Second + time.Second/10)

	pausedSched, err := s.svc.PauseSchedule(sched.ID)
	s.NoError(err)

	s.Equal(pausedSched.Status, model.ScheduleStatusPaused)
	pausedSched.Status = model.ScheduleStatusActive
	s.Equal(sched, pausedSched)

	calls := s.webhookHandlerCalls.Load()
	s.GreaterOrEqual(calls, int32(1))

	resumedSched, err := s.svc.ResumeSchedule(sched.ID)
	s.NoError(err)
	s.Equal(sched, resumedSched)

	time.Sleep(time.Second + time.Second/10)

	s.GreaterOrEqual(s.webhookHandlerCalls.Load(), calls)
}

type mockCronRepo struct {
	mtx    sync.Mutex
	nextID int64
	m      map[int64]*model.CronSchedule
}

func (s *mockCronRepo) Get(id int64) (*model.CronSchedule, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	sched, has := s.m[id]
	if !has {
		return nil, store.ErrScheduleNotExist
	}

	var copy model.CronSchedule = *sched
	return &copy, nil
}

func (s *mockCronRepo) Save(sched *model.CronSchedule) (int64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, has := s.m[sched.ID]; has {
		return sched.ID, nil
	}

	sched.ID = s.nextID
	s.nextID++

	s.m[sched.ID] = sched

	return sched.ID, nil
}

func (s *mockCronRepo) Delete(id int64) error {
	delete(s.m, id)

	return nil
}

func (s *mockCronRepo) Iter(iterFunc func(*model.CronSchedule) error) error {
	for _, cron := range s.m {
		if err := iterFunc(cron); err != nil {
			return err
		}
	}
	return nil
}

type mockStore struct {
	cronRepo *mockCronRepo
}

func (s *mockStore) CronScheduleRepository() store.CronScheduleRepository {
	if s.cronRepo == nil {
		s.cronRepo = &mockCronRepo{m: make(map[int64]*model.CronSchedule)}
	}
	return s.cronRepo
}

func (s *mockStore) HistoryRepository() store.CronHistoryRepository {
	return &mockHistoryRepo{}
}

type mockHistoryRepo struct {
	store.CronHistoryRepository
}

func (r *mockHistoryRepo) Insert(status *model.CronStatus, maxSamplesPerCron int) error {
	return nil
}
