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
	IsRecurring    bool
	Metadata       Metadata
	CreatedAt      time.Time
	NextScheduleAt time.Time
	StartAt        time.Time
	EndAt          time.Time
	Failures       int
}
