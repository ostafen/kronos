package db

import (
	"database/sql"

	"github.com/ostafen/kronos/internal/config"
)

func sqliteFilename(filename string) string {
	if filename == "" {
		return "kronos.sqlite"
	}
	return filename
}

func Open(conf config.Store) (*sql.DB, error) {
	switch conf.Driver {
	case "sqlite", "sqlite3":
		return sql.Open("sqlite3", sqliteFilename(conf.Host))
	}
	return nil, nil
}
