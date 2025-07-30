package db

import (
	"context"
	"database/sql"
	"ed-tracker/internal/config"
	"fmt"
	_ "modernc.org/sqlite"
)

type DBFactory struct {
	Path string
}

func (f *DBFactory) Connect(ctx context.Context) (*Queries, *sql.DB, error) {
	cfg := config.Load()

	if f.Path == "" {
		f.Path = cfg.DbFile
	}

	connStr := fmt.Sprintf("file:%s?cache=shared&_fk=1", f.Path)

	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("opening db: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("pinging db: %w", err)
	}

	queries := New(db)
	return queries, db, nil
}
