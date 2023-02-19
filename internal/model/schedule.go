package model

import (
	"fmt"
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/uuid"
	"github.com/ostafen/kronos/internal/cron"
	"github.com/ostafen/kronos/internal/db/dto"
	"github.com/sirupsen/logrus"
)

const (
	NotStartedStatus = "not_started"
	ActiveStatus     = "active"
	PausedStatus     = "paused"
	ElapsedStatus    = "elapsed"
)

type ScheduleRegisterInput struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	CronExpr    string            `json:"cronExpr"`
	Email       string            `json:"email" validate:"required"`
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

func (input *ScheduleRegisterInput) ToSched() (*dto.Schedule, error) {
	if err := validate(input); err != nil {
		return nil, err
	}

	nextSchedule := input.RunAt
	if input.Recurring() {
		start := input.StartAt
		if now := time.Now(); now.After(start) {
			start = now
		}
		nextSchedule = cron.NextTickAfter(input.CronExpr, start, true)
		logrus.Info(nextSchedule)
	}

	return &dto.Schedule{
		ID:             uuid.NewString(),
		Active:         true,
		Title:          input.Title,
		Description:    input.Description,
		CronExpr:       input.CronExpr,
		IsRecurring:    input.Recurring(),
		Email:          input.Email,
		URL:            input.URL,
		Metadata:       input.Metadata,
		NextScheduleAt: nextSchedule,
		StartAt:        input.StartAt,
		EndAt:          input.EndAt,
		CreatedAt:      time.Now(),
	}, nil
}

type Schedule struct {
	ID             string            `json:"id"`
	Title          string            `json:"title"`
	Status         string            `json:"status"`
	Description    string            `json:"description"`
	CronExpr       string            `json:"cronExpr"`
	Email          string            `json:"email"`
	URL            string            `json:"url"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      time.Time         `json:"createdAt"`
	NextScheduleAt time.Time         `json:"nextScheduleAt"`
	IsRecurring    bool              `json:"isRecurring"`
	RunAt          time.Time         `json:"runAt"`
	StartAt        time.Time         `json:"startAt"`
	EndAt          time.Time         `json:"endAt"`
}
