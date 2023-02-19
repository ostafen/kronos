package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ostafen/kronos/internal/cron"
	"github.com/ostafen/kronos/internal/db"
	"github.com/ostafen/kronos/internal/db/dto"
	"github.com/ostafen/kronos/internal/model"

	log "github.com/sirupsen/logrus"
)

type ScheduleService interface {
	RegisterSchedule(sched *model.ScheduleRegisterInput) (*model.Schedule, error)
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

func (s *schedService) RegisterSchedule(input *model.ScheduleRegisterInput) (*model.Schedule, error) {
	dtoSched, err := input.ToSched()
	if err != nil {
		return nil, err
	}

	_, err = s.schedRepo.Insert(dtoSched)
	if err != nil {
		return nil, err
	}

	sched := fromDTOSchedule(dtoSched)

	s.runCallbacks(sched, s.onScheduleRegistered)

	return sched, nil
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

	succeededSchedules, failedSchedules := s.nofifyAll(ctx, schedules)
	pausedSchedules, failedSchedules, err := s.pauseOrIncrementFailures(tx, failedSchedules)
	if err != nil {
		return nil, err
	}

	activeSchedules := concat(succeededSchedules, failedSchedules)
	if err := s.planNextSchedule(tx, activeSchedules); err != nil {
		return nil, err
	}

	if err := s.disableElapsedSchedules(tx, activeSchedules); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err == nil {
		s.notifyPausedSchedules(pausedSchedules...)
	}

	now := time.Now()
	return &now, err
}

func concat[T any](a []T, b []T) []T {
	res := make([]T, 0, len(a)+len(b))
	res = append(res, a...)
	res = append(res, b...)
	return res
}
func (s *schedService) disableElapsedSchedules(tx *sql.Tx, schedules []*dto.Schedule) error {
	now := time.Now()
	for _, sched := range schedules {
		if sched.EndAt.Before(now) {
			if err := s.schedRepo.UpdateActive(tx, sched.ID, false); err != nil {
				return err
			}
		}
	}
	return nil
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
		log.WithField("scheduleId", sched.ID).WithField("tick", sched.NextScheduleAt).Info("sendingNotification")
		err := s.notificationSvc.Send(ctx, sched.URL, sched)
		ch <- result{err: err, idx: i}
	}()
}

const (
	MaxFailures = 10
)

func (s *schedService) pauseOrIncrementFailures(tx *sql.Tx, schedules []*dto.Schedule) ([]*dto.Schedule, []*dto.Schedule, error) {
	paused := make([]*dto.Schedule, 0, len(schedules))
	failed := make([]*dto.Schedule, 0, len(schedules))

	for _, sched := range failed {
		log.WithField("scheduledId", sched.ID).WithField("failures", sched.Failures).Warn("schedule notification failed")

		newFailures := sched.Failures + 1
		if newFailures >= MaxFailures {
			log.WithField("scheduledId", sched.ID).Warn("pausing schedule")

			if err := s.schedRepo.UpdateActive(tx, sched.ID, false); err != nil {
				return nil, nil, err
			}
			paused = append(paused, sched)
		} else {
			failed = append(failed, sched)
			if err := s.schedRepo.UpdateFailures(tx, sched.ID, newFailures); err != nil {
				return nil, nil, err
			}
		}
	}
	return paused, failed, nil
}

func (s *schedService) planNextSchedule(tx *sql.Tx, schedules []*dto.Schedule) error {
	for _, sched := range schedules {
		nextScheduleTime := cron.NextTickAfter(sched.CronExpr, sched.NextScheduleAt, false)

		log.WithField("scheduleId", sched.ID).
			WithField("nextScheduleTime", nextScheduleTime).
			Info("planning next schedule")

		if err := s.schedRepo.UpdateScheduleTime(tx, sched.ID, nextScheduleTime); err != nil {
			return err
		}
	}
	return nil
}

func getStatus(isActive bool, startTime, endTime time.Time) string {
	now := time.Now()
	if now.Before(startTime) {
		return model.NotStartedStatus
	}

	if now.After(endTime) {
		return model.ElapsedStatus
	}

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
	var runAt time.Time
	if !sched.IsRecurring {
		runAt = sched.NextScheduleAt
	}

	return &model.Schedule{
		ID:             sched.ID,
		Title:          sched.Title,
		Description:    sched.Description,
		Status:         getStatus(sched.Active, sched.StartAt, sched.EndAt),
		Email:          sched.Email,
		URL:            sched.URL,
		CronExpr:       sched.CronExpr,
		Metadata:       sched.Metadata,
		IsRecurring:    sched.IsRecurring,
		StartAt:        sched.StartAt,
		RunAt:          runAt,
		EndAt:          sched.EndAt,
		CreatedAt:      sched.CreatedAt,
		NextScheduleAt: sched.NextScheduleAt,
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

	err = s.notificationSvc.Send(context.Background(), sched.URL, sched)
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

	nextTick := cron.NextTickAfter(dtoSched.CronExpr, time.Now(), false)
	log.WithField("scheduleId", dtoSched.ID).
		WithField("nextScheduleAt", nextTick).
		Info("resuming schedule")

	if err := s.schedRepo.UpdateFailures(tx, dtoSched.ID, 0); err != nil {
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
