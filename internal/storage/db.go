package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type DB struct{ SQL *sql.DB }

func Open(ctx context.Context, path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(8)
	if _, err = db.ExecContext(ctx, "PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("pragma: %w", err)
	}
	return &DB{SQL: db}, nil
}
func (d *DB) Close() error                   { return d.SQL.Close() }
func (d *DB) Ping(ctx context.Context) error { return d.SQL.PingContext(ctx) }
func (d *DB) Tx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
