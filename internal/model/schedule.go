package model

import (
	"fmt"
	"time"

	"github.com/ostafen/kronos/internal/cron"
)

type ScheduleStatus string

const (
	ScheduleStatusNotStarted ScheduleStatus = "not_started"
	ScheduleStatusActive     ScheduleStatus = "active"
	ScheduleStatusPaused     ScheduleStatus = "paused"
	ScheduleStatusExpired    ScheduleStatus = "expired"
)

type ScheduleRegisterInput struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	CronExpr    string            `json:"cronExpr"`
	URL         string            `json:"url" validate:"required"`
	IsRecurring *bool             `json:"isRecurring" validate:"required"`
	RunAt       time.Time         `json:"runAt"`
	StartAt     time.Time         `json:"startAt"`
	EndAt       time.Time         `json:"endAt"`
	Metadata    map[string]string `json:"metadata"`
}

func (input *ScheduleRegisterInput) Recurring() bool {
	return *input.IsRecurring
}

func validate(input *ScheduleRegisterInput) error {
	if input.Recurring() {
		if !cron.IsValid(input.CronExpr) {
			return fmt.Errorf("invalid cronExpr %s", input.CronExpr)
		}

		if !input.EndAt.IsZero() && input.EndAt.Before(time.Now()) {
			return fmt.Errorf(`"endAt" must be a valide date in the future`)
		}

		if input.EndAt.Before(input.StartAt) {
			return fmt.Errorf(`"endAt" must be greater than or equal to "startAt"`)
		}
	} else {
		if input.RunAt.IsZero() {
			return fmt.Errorf(`"runAt" must be set for non recurring schedule`)
		}

		if time.Now().After(input.RunAt) {
			return fmt.Errorf(`"runAt" must be a valide date in the future`)
		}

		if !input.StartAt.IsZero() || !input.EndAt.IsZero() {
			return fmt.Errorf(`"startAt"/"endAt" should not be set together with "runAt"`)
		}
	}
	return nil
}

var maxTime = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

func (input *ScheduleRegisterInput) ToSched() (*CronSchedule, error) {
	if err := validate(input); err != nil {
		return nil, err
	}

	startAt, endAt := input.StartAt, input.EndAt
	if !input.RunAt.IsZero() {
		startAt = input.RunAt
		endAt = input.RunAt
	} else if input.EndAt.IsZero() {
		endAt = maxTime
	}

	return &CronSchedule{
		ID:          -1,
		Status:      ScheduleStatusActive,
		Title:       input.Title,
		Description: input.Description,
		CronExpr:    input.CronExpr,
		IsRecurring: input.Recurring(),
		URL:         input.URL,
		Metadata:    input.Metadata,
		RunAt:       input.RunAt,
		StartAt:     startAt,
		EndAt:       endAt,
		CreatedAt:   time.Now(),
	}, nil
}

type CronSchedule struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Status      ScheduleStatus    `json:"status"`
	Description string            `json:"description"`
	CronExpr    string            `json:"cronExpr"`
	URL         string            `json:"url"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"createdAt"`
	IsRecurring bool              `json:"isRecurring"`
	RunAt       time.Time         `json:"runAt,omitempty"`
	StartAt     time.Time         `json:"startAt"`
	EndAt       time.Time         `json:"endAt"`
	Failures    int               `json:"-"`
}

func (s *CronSchedule) nextTick(start time.Time) time.Time {
	if !s.IsRecurring {
		return s.RunAt
	}
	return cron.Next(s.CronExpr, start)
}

func (s *CronSchedule) Expired() bool {
	return !s.EndAt.After(time.Now())
}

func (s *CronSchedule) NextTick() time.Time {
	now := time.Now()
	if s.StartAt.After(now) {
		return s.nextTick(s.StartAt)
	}
	return s.nextTick(now)
}

func (s *CronSchedule) IsActive() bool {
	return s.Status == ScheduleStatusActive
}

type CronStatus struct {
	CronID     int64         `json:"cronId"`
	At         time.Time     `json:"at"`
	StatusCode int           `json:"statusCode"`
	Duration   time.Duration `json:"duration"`
}
