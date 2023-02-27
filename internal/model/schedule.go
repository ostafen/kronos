package model

import (
	"fmt"
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
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
	StartAt     time.Time         `json:"startAt" validate:"required"`
	EndAt       time.Time         `json:"endAt" validate:"required"`
	Metadata    map[string]string `json:"metadata"`
}

func (input *ScheduleRegisterInput) Recurring() bool {
	return *input.IsRecurring
}

func validate(input *ScheduleRegisterInput) error {
	if *input.IsRecurring {
		gron := gronx.New()
		if !gron.IsValid(input.CronExpr) {
			return fmt.Errorf("invalid cronExpr %s", input.CronExpr)
		}

		if input.StartAt.IsZero() || input.EndAt.IsZero() {
			return fmt.Errorf(`both "startAt" and "endAt" must be specified for recurring schedule`)
		}

		if input.EndAt.Before(time.Now()) {
			return fmt.Errorf(`"endAt" must be a valide date in the future`)
		}

		if input.EndAt.Before(input.StartAt) {
			return fmt.Errorf(`"endAt" must be greater than or equal to "startAt"`)
		}
	}

	if !input.Recurring() {
		if input.RunAt.IsZero() {
			return fmt.Errorf(`"runAt" must be set for non recurring schedule`)
		}

		if time.Now().After(input.RunAt) {
			return fmt.Errorf(`"runAt" must be a valide date in the future`)
		}

		if !input.StartAt.Equal(input.RunAt) || !input.EndAt.Equal(input.RunAt) {
			return fmt.Errorf(`"startAt" and "endAt" must both be equal to "runAt"`)
		}
	}
	return nil
}

func (input *ScheduleRegisterInput) ToSched() (*Schedule, error) {
	if err := validate(input); err != nil {
		return nil, err
	}

	return &Schedule{
		ID:          uuid.NewString(),
		Status:      ScheduleStatusActive,
		Title:       input.Title,
		Description: input.Description,
		CronExpr:    input.CronExpr,
		IsRecurring: input.Recurring(),
		URL:         input.URL,
		Metadata:    input.Metadata,
		RunAt:       input.RunAt,
		StartAt:     input.StartAt,
		EndAt:       input.EndAt,
		CreatedAt:   time.Now(),
	}, nil
}

type Schedule struct {
	ID          string            `json:"id"`
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

func (s *Schedule) nextTick(start time.Time, includeStart bool) time.Time {
	if !s.IsRecurring {
		return s.RunAt
	}
	return cron.NextTickAfter(s.CronExpr, start, includeStart)
}

func (s *Schedule) FirstTick() time.Time {
	return s.nextTick(s.StartAt, true)
}

func (s *Schedule) Expired() bool {
	return !s.EndAt.After(time.Now())
}

func (s *Schedule) NextTickAt() time.Time {
	now := time.Now()
	if now.Before(s.FirstTick()) {
		return s.FirstTick()
	}
	return s.nextTick(time.Now(), false)
}

func (s *Schedule) IsActive() bool {
	return s.Status == ScheduleStatusActive
}
