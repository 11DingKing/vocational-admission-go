package repository

import (
	"context"
	"database/sql"
	"time"
)

type HealthRepo struct{ DB *sql.DB }

func (h HealthRepo) Probe(ctx context.Context) error {
	var one int
	return h.DB.QueryRowContext(ctx, "SELECT 1").Scan(&one)
}
func (h HealthRepo) OldestAudit(ctx context.Context) (time.Time, error) {
	var ts sql.NullString
	e := h.DB.QueryRowContext(ctx, "SELECT MIN(created_at) FROM audit_events").Scan(&ts)
	if e != nil {
		return time.Time{}, e
	}
	if !ts.Valid {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, ts.String)
}
func (h HealthRepo) TableCount(ctx context.Context) (int, error) {
	var n int
	e := h.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='table'").Scan(&n)
	return n, e
}
