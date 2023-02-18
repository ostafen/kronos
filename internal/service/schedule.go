package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ostafen/kronos/internal/db"
	"github.com/ostafen/kronos/internal/db/dto"
	"github.com/ostafen/kronos/internal/model"

	log "github.com/sirupsen/logrus"
)

type ScheduleService interface {
	RegisterSchedule(sched *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	DeleteSchedule(id string) error
	ListSchedules(offset, limit int) ([]*model.Schedule, error)

	PauseSchedule(id string) (*model.Schedule, error)
	ResumeSchedule(id string) (*model.Schedule, error)
	TriggerSchedule(id string) (*model.Schedule, error)

	NotifySchedules() (*time.Time, error)

	OnScheduleResumed(cbk ...ScheduleCallback)
	OnScheduleRegistered(cbk ...ScheduleCallback)
	OnSchedulePaused(cbk ...ScheduleCallback)
}

func NewScheduleService(db *sql.DB, schedRepo db.ScheduleRepo, svc NotificationService) ScheduleService {
	return &schedService{
		dbConn:          db,
		schedRepo:       schedRepo,
		notificationSvc: svc,
	}
}

type schedService struct {
	dbConn    *sql.DB
	schedRepo db.ScheduleRepo

	notificationSvc NotificationService

	onScheduleRegistered []ScheduleCallback
	onScheduleResumed    []ScheduleCallback
	onSchedulePaused     []ScheduleCallback
}

func (s *schedService) RegisterSchedule(sched *model.Schedule) error {
	dtoSched := &dto.Schedule{
		ID:             uuid.NewString(),
		Active:         true,
		Title:          sched.Title,
		Description:    sched.Description,
		Email:          sched.Email,
		URL:            sched.URL,
		CronExpr:       sched.CronExpr,
		Metadata:       sched.Metadata,
		NextScheduleAt: sched.FirstSchedule(),
		CreatedAt:      time.Now().UTC(),
	}

	id, err := s.schedRepo.Insert(dtoSched)
	if err != nil {
		return err
	}

	sched.ID = id
	sched.CreatedAt = dtoSched.CreatedAt
	sched.Status = getStatus(dtoSched.Active)

	s.runCallbacks(sched, s.onScheduleRegistered)

	return nil
}

const (
	MaxConcurrentNotifications = 100
	MaxRequestDuration         = time.Second * 5
)

func (s *schedService) NotifySchedules() (*time.Time, error) {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	schedules, err := s.schedRepo.PickPending(tx, MaxConcurrentNotifications)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		time, err := s.schedRepo.NextScheduleTime(tx)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return time, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), MaxRequestDuration)
	defer cancel()

	succeeded, failed := s.nofifyAll(ctx, schedules)
	paused, err := s.planNextSchedule(tx, succeeded, failed)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err == nil {
		s.notifyPausedSchedules(paused...)
	}

	now := time.Now()
	return &now, err
}

func (s *schedService) nofifyAll(ctx context.Context, schedules []*dto.Schedule) ([]*dto.Schedule, []*dto.Schedule) {
	ch := make(chan result, len(schedules))
	for i, sched := range schedules {
		s.notify(ctx, fromDTOSchedule(sched), i, ch)
	}

	succeeded := make([]*dto.Schedule, 0, len(schedules))
	failed := make([]*dto.Schedule, 0, len(schedules))
	for {
		select {
		case <-ctx.Done():
			return succeeded, failed
		case res := <-ch:
			sched := schedules[res.idx]
			if res.err != nil {
				failed = append(failed, sched)
			} else {
				succeeded = append(succeeded, sched)
			}
		}

		if len(failed)+len(succeeded) == len(schedules) {
			return succeeded, failed
		}
	}
}

type result struct {
	err error
	idx int
}

func (s *schedService) notify(ctx context.Context, sched *model.Schedule, i int, ch chan result) {
	go func() {
		err := s.notificationSvc.Send(ctx, sched.URL, sched)
		ch <- result{err: err, idx: i}
	}()
}

const (
	MaxFailures = 1
)

func (s *schedService) planNextSchedule(tx *sql.Tx, succeeded, failed []*dto.Schedule) ([]*dto.Schedule, error) {
	schedules := make([]*dto.Schedule, 0, len(succeeded)+len(failed))

	schedules = append(schedules, succeeded...)
	schedules = append(schedules, failed...)

	paused := make([]*dto.Schedule, 0, len(failed))
	for i, sched := range schedules {
		failed := i >= len(succeeded)
		if failed && sched.Failures+1 >= MaxFailures {
			if err := s.schedRepo.UpdateActive(tx, sched.ID, false); err != nil {
				return nil, err
			}
			paused = append(paused, sched)
		}

		nextScheduleTime := model.NextScheduleAfter(sched.CronExpr, sched.NextScheduleAt)

		log.WithField("scheduleId", sched.ID).
			WithField("nextScheduleTime", nextScheduleTime).
			Info("planning next schedule")

		failures := 0
		if failed {
			failures = sched.Failures + 1
		}

		if err := s.schedRepo.UpdateScheduleTimeAndFailures(tx, sched.ID, nextScheduleTime, failures); err != nil {
			return nil, err
		}
	}

	return paused, nil
}

func getStatus(isActive bool) string {
	if isActive {
		return model.ActiveStatus
	}
	return model.PausedStatus
}

func (s *schedService) GetSchedule(id string) (*model.Schedule, error) {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	dtoSchedule, err := s.schedRepo.Get(tx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return fromDTOSchedule(dtoSchedule), nil
}

func (s *schedService) DeleteSchedule(id string) error {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.schedRepo.Delete(tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *schedService) ListSchedules(offset, limit int) ([]*model.Schedule, error) {
	dtoSchedules, err := s.schedRepo.List()
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	schedules := make([]*model.Schedule, 0, len(dtoSchedules))
	for _, sched := range dtoSchedules {
		s := fromDTOSchedule(sched)
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func fromDTOSchedule(sched *dto.Schedule) *model.Schedule {
	return &model.Schedule{
		ID:             sched.ID,
		Title:          sched.Title,
		Description:    sched.Description,
		Status:         getStatus(sched.Active),
		Email:          sched.Email,
		URL:            sched.URL,
		CronExpr:       sched.CronExpr,
		Metadata:       sched.Metadata,
		CreatedAt:      sched.CreatedAt,
		NextScheduleAt: sched.NextScheduleAt,
		Failures:       sched.Failures,
	}
}

func (s *schedService) pauseSchedule(tx *sql.Tx, id string) error {
	return s.schedRepo.UpdateActive(tx, id, false)
}

func (s *schedService) resumeSchedule(tx *sql.Tx, id string) error {
	return s.schedRepo.UpdateActive(tx, id, true)
}

func (s *schedService) notifyPausedSchedules(schedules ...*dto.Schedule) {
	for _, sched := range schedules {
		s.runCallbacks(fromDTOSchedule(sched), s.onSchedulePaused)
	}
}

func (s *schedService) runCallbacks(sched *model.Schedule, callbacks []ScheduleCallback) {
	for _, cbk := range callbacks {
		cbk(sched)
	}
}

func (s *schedService) PauseSchedule(id string) (*model.Schedule, error) {
	log.WithField("scheduleId", id).Info("pausing schedule")

	tx, err := s.dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sched, err := s.schedRepo.Get(tx, id)
	if err != nil {
		return nil, err
	}

	if err := s.pauseSchedule(tx, sched.ID); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err == nil {
		s.notifyPausedSchedules(sched)
	}
	return fromDTOSchedule(sched), err
}

func (s *schedService) TriggerSchedule(id string) (*model.Schedule, error) {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	sched, err := s.schedRepo.Get(tx, id)
	if err != nil {
		return nil, err
	}

	err = s.notificationSvc.Send(context.Background(), sched.URL, &dto.ScheduleNotification{ScheduleID: sched.ID})
	return fromDTOSchedule(sched), err
}

func (s *schedService) ResumeSchedule(id string) (*model.Schedule, error) {
	tx, err := s.dbConn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	dtoSched, err := s.schedRepo.Get(tx, id)
	if err != nil {
		return nil, err
	}

	if err := s.resumeSchedule(tx, dtoSched.ID); err != nil {
		return nil, err
	}

	nextTick := model.NextScheduleAfter(dtoSched.CronExpr, time.Now())
	log.WithField("scheduleId", dtoSched.ID).
		WithField("nextScheduleAt", nextTick).
		Info("resuming schedule")

	if err := s.schedRepo.UpdateScheduleTimeAndFailures(tx, dtoSched.ID, nextTick, 0); err != nil {
		return nil, err
	}

	sched := fromDTOSchedule(dtoSched)

	err = tx.Commit()
	if err == nil {
		s.runCallbacks(sched, s.onScheduleResumed)
	}
	return sched, nil
}

type ScheduleCallback func(*model.Schedule)

func (s *schedService) OnScheduleResumed(callbacks ...ScheduleCallback) {
	s.onScheduleResumed = append(s.onScheduleResumed, callbacks...)
}

func (s *schedService) OnScheduleRegistered(callbacks ...ScheduleCallback) {
	s.onScheduleRegistered = append(s.onScheduleRegistered, callbacks...)
}

func (s *schedService) OnSchedulePaused(callbacks ...ScheduleCallback) {
	s.onSchedulePaused = append(s.onSchedulePaused, callbacks...)
}
