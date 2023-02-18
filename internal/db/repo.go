package db

import (
	"database/sql"
	"time"

	"github.com/ostafen/kronos/internal/db/dto"
)

type ScheduleRepo interface {
	List() ([]*dto.Schedule, error)
	Get(tx *sql.Tx, id string) (*dto.Schedule, error)
	Insert(sched *dto.Schedule) (string, error)
	Delete(tx *sql.Tx, id string) error
	PickPending(tx *sql.Tx, limit int) ([]*dto.Schedule, error)
	NextScheduleTime(tx *sql.Tx) (*time.Time, error)
	UpdateScheduleTimeAndFailures(tx *sql.Tx, id string, schedTime time.Time, failures int) error
	UpdateActive(tx *sql.Tx, id string, active bool) error
}

func GetScheduleRepo(dbConn *sql.DB) ScheduleRepo {
	return &sqliteJobRepo{
		DB: dbConn,
	}
}
