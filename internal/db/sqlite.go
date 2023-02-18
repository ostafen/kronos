package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ostafen/kronos/internal/db/dto"
)

type sqliteJobRepo struct {
	*sql.DB
}

func placeholdersOffset(n, offset int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("$%d", i+offset)
	}
	return s
}

func placeholders(n int) string {
	return placeholdersOffset(n, 1)
}

func (repo *sqliteJobRepo) Insert(sched *dto.Schedule) (string, error) {
	insertStmt := fmt.Sprintf(`INSERT INTO %s(%s) VALUES (%s) RETURNING %s`, schedTableName, strings.Join(schedTableCols, ","), placeholders(len(schedTableCols)), schedTableIdCol)

	row := repo.QueryRow(insertStmt,
		&sched.ID,
		&sched.Title,
		&sched.Active,
		&sched.Description,
		&sched.URL,
		&sched.CronExpr,
		&sched.Email,
		&sched.CreatedAt,
		&sched.NextScheduleAt,
		&sched.Metadata,
		&sched.Failures,
	)

	var jobID string
	if err := row.Scan(&jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

func (repo *sqliteJobRepo) Get(tx *sql.Tx, id string) (*dto.Schedule, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1`, strings.Join(schedTableCols, ","), schedTableName, schedTableIdCol)
	rows := tx.QueryRow(query, id)

	return ScanSchedule(rows)
}

func (repo *sqliteJobRepo) List() ([]*dto.Schedule, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(schedTableCols, ","), schedTableName)
	rows, err := repo.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]*dto.Schedule, 0)
	for rows.Next() {
		sched, err := ScanSchedule(rows)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, sched)
	}
	return schedules, nil
}

func ScanSchedule[R Row](row R) (*dto.Schedule, error) {
	var sched dto.Schedule

	err := row.Scan(
		&sched.ID,
		&sched.Title,
		&sched.Active,
		&sched.Description,
		&sched.URL,
		&sched.CronExpr,
		&sched.Email,
		&sched.CreatedAt,
		&sched.NextScheduleAt,
		&sched.Metadata,
		&sched.Failures,
	)
	return &sched, err
}

func (repo *sqliteJobRepo) PickPending(tx *sql.Tx, limit int) ([]*dto.Schedule, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1 AND %s <= $2 ORDER BY %s ASC LIMIT $3`,
		strings.Join(schedTableCols, ","),
		schedTableName,
		schedTableActiveCol,
		schedTableNextScheduleAtCol,
		schedTableNextScheduleAtCol,
	)

	rows, err := tx.Query(query, true, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scheds := make([]*dto.Schedule, 0)
	for rows.Next() {
		s, err := ScanSchedule(rows)
		if err != nil {
			return nil, err
		}
		scheds = append(scheds, s)
	}
	return scheds, nil
}

func (repo *sqliteJobRepo) NextScheduleTime(tx *sql.Tx) (*time.Time, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = $1 ORDER BY %s ASC LIMIT 1`, schedTableNextScheduleAtCol, schedTableName, schedTableActiveCol, schedTableNextScheduleAtCol)

	var nextTime *time.Time
	row := tx.QueryRow(query, true)
	err := row.Scan(&nextTime)

	return nextTime, err
}

func (repo *sqliteJobRepo) Delete(tx *sql.Tx, id string) error {
	stmt := fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, schedTableName, schedTableIdCol)

	_, err := tx.Exec(stmt, id)
	return err
}

func (repo *sqliteJobRepo) UpdateScheduleTimeAndFailures(tx *sql.Tx, id string, schedTime time.Time, failures int) error {
	stmt := fmt.Sprintf(`UPDATE %s SET %s = $1, %s = $2 WHERE %s = $3`,
		schedTableName,
		schedTableNextScheduleAtCol,
		schedTableFailuresCol,
		schedTableIdCol,
	)

	_, err := tx.Exec(stmt, schedTime, failures, id)
	return err
}

func (repo *sqliteJobRepo) UpdateActive(tx *sql.Tx, id string, active bool) error {
	stmt := fmt.Sprintf(`UPDATE %s SET %s = $1 WHERE %s = $2`,
		schedTableName,
		schedTableActiveCol,
		schedTableIdCol,
	)

	_, err := tx.Exec(stmt, active, id)
	return err
}
