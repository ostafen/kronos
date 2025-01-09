package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ostafen/kronos/internal/model"
)

var ErrScheduleNotExist = errors.New("schedule does not exist")

type Store interface {
	CronScheduleRepository() CronScheduleRepository
	HistoryRepository() CronHistoryRepository
}

type CronScheduleRepository interface {
	Get(id int64) (*model.CronSchedule, error)
	Save(sched *model.CronSchedule) (int64, error)
	Delete(id int64) error
	Iter(iterFunc func(cron *model.CronSchedule) error) error
}

type CronHistoryRepository interface {
	Insert(status *model.CronStatus, maxSamplesPerCron int) error
	GetHistory(n int) ([]*model.CronStatus, error)
	GetCronHistory(cronID int64, n int) ([]*model.CronStatus, error)
}

var (
	cronSchedulesCols = []string{
		"id",
		"title",
		"status",
		"description",
		"cron_expr",
		"url",
		"metadata",
		"created_at",
		"is_recurring",
		"run_at",
		"start_at",
		"end_at",
	}

	cronStatusCols = []string{
		"cron_id",
		"at",
		"status_code",
		"duration",
	}
)

var (
	cronSchedulesValues []string
	cronStatusValues    []string
)

func init() {
	cronSchedulesValues = make([]string, len(cronSchedulesCols))
	for i := range cronSchedulesCols {
		cronSchedulesValues[i] = fmt.Sprintf("$%d", i)
	}

	cronStatusValues = make([]string, len(cronStatusCols))
	for i := range cronStatusCols {
		cronStatusValues[i] = fmt.Sprintf("$%d", i)
	}
}

type sqlStore struct {
	db *sql.DB
}

func (s *sqlStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS cron_schedules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR NOT NULL,
			status VARCHAR NOT NULL,
			description VARCHAR NOT NULL,
			cron_expr VARCHAR NOT NULL,
			url VARCHAR NOT NULL,
			metadata VARCHAR,
			created_at TIMESTAMP NOT NULL,
			is_recurring BOOLEAN NOT NULL,
			run_at TIMESTAMP,
			start_at TIMESTAMP NOT NULL,
			end_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS cron_status (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cron_id INTEGER,
			at TIMESTAMP,
			status_code INTEGER,
			duration INTEGER
		);

		CREATE INDEX IF NOT EXISTS at_index ON cron_status(at);
	`)
	return err
}

func New(path string) (Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	s := &sqlStore{db: db}
	return s, s.migrate()
}

func (s *sqlStore) CronScheduleRepository() CronScheduleRepository {
	return &cronScheduleRepo{s.db}
}

func (s *sqlStore) HistoryRepository() CronHistoryRepository {
	return &statusRepo{db: s.db}
}

type cronScheduleRepo struct {
	db *sql.DB
}

func (s *cronScheduleRepo) Get(id int64) (*model.CronSchedule, error) {
	row := s.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM cron_schedules WHERE id = $1", strings.Join(cronSchedulesCols, ",")),
		id,
	)
	return scanCron(row)
}

func (s *cronScheduleRepo) Save(cron *model.CronSchedule) (int64, error) {
	metadata, err := json.Marshal(cron.Metadata)
	if err != nil {
		return -1, err
	}

	values := []any{
		cron.ID,
		cron.Title,
		cron.Status,
		cron.Description,
		cron.CronExpr,
		cron.URL,
		metadata,
		cron.CreatedAt,
		cron.IsRecurring,
		cron.RunAt,
		cron.StartAt,
		cron.EndAt,
	}

	cols := cronSchedulesCols
	placeHolders := cronSchedulesValues
	if cron.ID <= 0 {
		cols = cols[1:]
		values = values[1:]
		placeHolders = placeHolders[:len(placeHolders)-1]
	}

	row := s.db.QueryRow(
		fmt.Sprintf(
			`INSERT INTO cron_schedules(%s) VALUES (%s)
			ON CONFLICT (id) DO UPDATE
			SET title = excluded.title, status = excluded.status, description = excluded.description,
				cron_expr = excluded.cron_expr, url = excluded.url, metadata = excluded.metadata,
				is_recurring = excluded.is_recurring, run_at = excluded.run_at, start_at = excluded.start_at,
				end_at = excluded.end_at
			RETURNING id;
			`,
			strings.Join(cols, ","),
			strings.Join(placeHolders, ","),
		),
		values...,
	)

	var id int64
	err = row.Scan(&id)
	return id, err
}

func (s *cronScheduleRepo) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM cron_schedules WHERE id = $1", id)
	return err
}

func (s *cronScheduleRepo) Iter(onCron func(cron *model.CronSchedule) error) error {
	rows, err := s.db.Query(
		fmt.Sprintf("SELECT %s FROM cron_schedules", strings.Join(cronSchedulesCols, ",")),
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		cron, err := scanCron(rows)
		if err != nil {
			return err
		}

		if err := onCron(cron); err != nil {
			return err
		}
	}
	return nil
}

func scanCron[T interface{ Scan(...any) error }](row T) (*model.CronSchedule, error) {
	var cron model.CronSchedule
	var metadata string

	err := row.Scan(
		&cron.ID,
		&cron.Title,
		&cron.Status,
		&cron.Description,
		&cron.CronExpr,
		&cron.URL,
		&metadata,
		&cron.CreatedAt,
		&cron.IsRecurring,
		&cron.RunAt,
		&cron.StartAt,
		&cron.EndAt,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(metadata), &cron.Metadata)
	return &cron, err
}

type statusRepo struct {
	db *sql.DB
}

func (r *statusRepo) Insert(cs *model.CronStatus, maxSamplesPerCron int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		fmt.Sprintf(
			`INSERT INTO cron_status(%s) VALUES(%s)`,
			strings.Join(cronStatusCols, ","),
			strings.Join(cronStatusValues, ","),
		),
		cs.CronID,
		cs.At,
		cs.StatusCode,
		cs.Duration,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM cron_status WHERE cron_id = $1 AND id NOT IN (SELECT id FROM cron_status WHERE cron_id = $1 ORDER BY at DESC LIMIT $2)`, cs.CronID, maxSamplesPerCron)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *statusRepo) GetCronHistory(cronID int64, n int) ([]*model.CronStatus, error) {
	statuses := make([]*model.CronStatus, 0, n)

	rows, err := r.db.Query(
		fmt.Sprintf("SELECT %s FROM cron_status WHERE cron_id = $1 ORDER BY at DESC LIMIT $2", strings.Join(cronStatusCols, ",")),
		cronID,
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s model.CronStatus
		err := rows.Scan(
			&s.CronID,
			&s.At,
			&s.StatusCode,
			&s.Duration,
		)
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, &s)
	}
	return statuses, nil
}

func (r *statusRepo) GetHistory(n int) ([]*model.CronStatus, error) {
	statuses := make([]*model.CronStatus, 0, n)

	rows, err := r.db.Query(
		fmt.Sprintf("SELECT %s FROM cron_status ORDER BY at DESC LIMIT $2", strings.Join(cronStatusCols, ",")),
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s model.CronStatus
		err := rows.Scan(
			&s.CronID,
			&s.At,
			&s.StatusCode,
			&s.Duration,
		)
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, &s)
	}
	return statuses, nil
}
