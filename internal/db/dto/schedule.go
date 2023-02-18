package dto

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Metadata map[string]string

func (meta *Metadata) Value() (driver.Value, error) {
	return json.Marshal(meta)
}

func (meta *Metadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &meta)
}

type Schedule struct {
	ID             string
	Active         bool
	Title          string
	Description    string
	CronExpr       string
	URL            string
	Email          string
	Metadata       Metadata
	NextScheduleAt time.Time
	CreatedAt      time.Time
	Failures       int
}

type ScheduleNotification struct {
	ScheduleID string `json:"scheduleId"`
}

type ScheduleStatus string

const (
	ScheduleStatusNotified ScheduleStatus = "NOTIFIED"
	ScheduleStatusFailed   ScheduleStatus = "FAILED"
)

type ScheduleHistoryItem struct {
	ScheduleID  string
	EndpointID  int
	ScheduledAt time.Time
	NotifiedAt  time.Time
	Status      ScheduleStatus
}
