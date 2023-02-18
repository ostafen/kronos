package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Migrator struct {
	dbConn *sql.DB
}

func NewMigrator(dbConn *sql.DB) *Migrator {
	return &Migrator{
		dbConn: dbConn,
	}
}

func (m *Migrator) Migrate() error {
	tx, err := m.dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := m.createSchedulesTable(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (m *Migrator) createSchedulesTable(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE IF NOT EXISTS schedules (
		id VARCHAR(100) NOT NULL PRIMARY KEY,
		title VARCHAR(100) NOT NULL UNIQUE,
		description VARCHAR(255) NULL DEFAULT NULL,
		active BOOL NOT NULL DEFAULT TRUE,
		url TEXT NOT NULL,
		cron_expr VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		next_schedule_at TIMESTAMP NOT NULL,
		failures INT NOT NULL DEFAULT 0,
		metadata BLOB NOT NULL
	)`)
	return err
}
