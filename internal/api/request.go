package api

import (
	"github.com/ostafen/kronos/internal/model"
)

type ScheduleCreateRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	CronExpr    string            `json:"cronExpr" validate:"required"`
	Email       string            `json:"email" validate:"required"`
	URL         string            `json:"url" validate:"required"`
	Metadata    map[string]string `json:"metadata"`
}

func (r *ScheduleCreateRequest) ToSched() *model.Schedule {
	return &model.Schedule{
		Title:       r.Name,
		Description: r.Description,
		CronExpr:    r.CronExpr,
		Email:       r.Email,
		URL:         r.URL,
		Metadata:    r.Metadata,
	}
}
