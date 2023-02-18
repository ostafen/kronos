package api

import "time"

type ScheduleJobResponse struct {
	ScheduleID   string    `json:"scheduleId"`
	NextSchedule time.Time `json:"nextSchedule"`
}
