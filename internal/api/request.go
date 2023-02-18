package api

import (
	"github.com/ostafen/kronos/internal/model"
)

type ScheduleCreateRequest struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description"`
	CronExpr    string            `json:"cronExpr" validate:"required"`
	Email       string            `json:"email" validate:"required"`
	URL         string            `json:"url" validate:"required"`
	Metadata    map[string]string `json:"metadata"`
}

func (r *ScheduleCreateRequest) ToSched() *model.Schedule {
	return &model.Schedule{
		Title:       r.Title,
		Description: r.Description,
		CronExpr:    r.CronExpr,
		Email:       r.Email,
		URL:         r.URL,
		Metadata:    r.Metadata,
	}
}
