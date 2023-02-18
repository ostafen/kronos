package model

import (
	"log"
	"time"

	"github.com/adhocore/gronx"
)

const (
	ActiveStatus = "active"
	PausedStatus = "paused"
)

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
	Failures       int               `json:"failures"`
}

func (s *Schedule) FirstSchedule() time.Time {
	t, err := gronx.NextTick(s.CronExpr, true)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

type Endpoint struct {
	URL      string `json:"url"`
	IsBroken bool   `json:"isBroken"`
}

func NextScheduleAfter(cronExpr string, time time.Time) time.Time {
	t, err := gronx.NextTickAfter(cronExpr, time, false)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
