package service

import (
	"context"
	"errors"
	"time"

	"github.com/ostafen/kronos/internal/model"
	"github.com/ostafen/kronos/internal/store"

	log "github.com/sirupsen/logrus"
)

type ScheduleService interface {
	RegisterSchedule(sched *model.ScheduleRegisterInput) (*model.Schedule, error)
	GetSchedule(id string) (*model.Schedule, error)
	DeleteSchedule(id string) error
	IterSchedules(onSchedule func(*model.Schedule) error) error

	OnTick(schedID string) time.Time
	PauseSchedule(id string) (*model.Schedule, error)
	ResumeSchedule(id string) (*model.Schedule, error)
	TriggerSchedule(id string) (*model.Schedule, error)

	OnScheduleResumed(cbk ...ScheduleCallback)
	OnScheduleRegistered(cbk ...ScheduleCallback)
	OnSchedulePaused(cbk ...ScheduleCallback)
	OnScheduleNotified(cbk ...ScheduleNotificationCallback)
}

func NewScheduleService(store store.ScheduleStore, svc NotificationService) ScheduleService {
	return &schedService{
		store:           store,
		notificationSvc: svc,
	}
}

type schedService struct {
	notificationSvc NotificationService

	store                store.ScheduleStore
	onScheduleRegistered []ScheduleCallback
	onScheduleResumed    []ScheduleCallback
	onSchedulePaused     []ScheduleCallback
	onScheduleNotified   []ScheduleNotificationCallback
}

func (s *schedService) RegisterSchedule(input *model.ScheduleRegisterInput) (*model.Schedule, error) {
	sched, err := input.ToSched()
	if err != nil {
		return nil, err
	}

	err = s.store.InsertOrUpdate(sched)
	if err == nil {
		s.runCallbacks(sched, s.onScheduleRegistered)
	}
	return sched, err
}

const (
	MaxRequestDuration = time.Second * 5
)

func (s *schedService) OnTick(schedID string) time.Time {
	sched, err := s.store.Get(schedID)
	if errors.Is(err, store.ErrNotPresent) {
		log.Errorf("no schedule with id %s", schedID)
		return time.Time{}
	}

	if err != nil {
		log.Fatal(err)
	}

	go s.sendWebhookNotification(sched)

	if sched.Expired() {
		return time.Time{}
	}
	return sched.NextTickAt()
}

func (s *schedService) sendWebhookNotification(sched *model.Schedule) error {
	log.WithField("scheduleId", sched.ID).
		WithField("url", sched.URL).
		Info("sendingNotification")

	ctx, cancel := context.WithTimeout(context.Background(), MaxRequestDuration)
	defer cancel()
	code, err := s.notificationSvc.Send(ctx, sched.URL, sched)
	if err == nil {
		s.runNotificationCallbacks(sched, code)
	}
	return err
}

func (s *schedService) runNotificationCallbacks(sched *model.Schedule, code int) {
	for _, cbk := range s.onScheduleNotified {
		cbk(sched, code)
	}
}

func (s *schedService) GetSchedule(id string) (*model.Schedule, error) {
	sched, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return sched, nil
}

func (s *schedService) DeleteSchedule(id string) error {
	return s.store.Delete(id)
}

func (s *schedService) IterSchedules(onSched func(*model.Schedule) error) error {
	return s.store.Iterate(onSched)
}

func (s *schedService) runCallbacks(sched *model.Schedule, callbacks []ScheduleCallback) {
	for _, cbk := range callbacks {
		cbk(sched)
	}
}

func (s *schedService) PauseSchedule(id string) (*model.Schedule, error) {
	log.WithField("scheduleId", id).Info("pausing schedule")

	sched, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	sched.Status = model.ScheduleStatusPaused
	if err := s.store.InsertOrUpdate(sched); err != nil {
		return nil, err
	}

	s.runCallbacks(sched, s.onSchedulePaused)
	return sched, nil
}

func (s *schedService) TriggerSchedule(id string) (*model.Schedule, error) {
	sched, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	err = s.sendWebhookNotification(sched)
	return sched, err
}

func (s *schedService) ResumeSchedule(id string) (*model.Schedule, error) {
	sched, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	sched.Status = model.ScheduleStatusActive
	if err := s.store.InsertOrUpdate(sched); err != nil {
		return nil, err
	}

	log.WithField("scheduleId", sched.ID).
		WithField("nextScheduleAt", sched.NextTickAt()).
		Info("resuming schedule")

	s.runCallbacks(sched, s.onScheduleResumed)
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

type ScheduleNotificationCallback func(*model.Schedule, int)

func (s *schedService) OnScheduleNotified(callbacks ...ScheduleNotificationCallback) {
	s.onScheduleNotified = append(s.onScheduleNotified, callbacks...)
}
