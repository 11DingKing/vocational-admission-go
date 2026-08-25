package repository

import (
	"context"
	"database/sql"
	"time"
)

func CleanupExpired(ctx context.Context, db *sql.DB, now time.Time) (int64, error) {
	res, e := db.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at<?", now.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
func PendingJobs(ctx context.Context, db *sql.DB) (int, error) {
	var n int
	e := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM jobs WHERE completed_at IS NULL").Scan(&n)
	return n, e
}
func MarkStaleJobs(ctx context.Context, db *sql.DB, before time.Time) (int64, error) {
	res, e := db.ExecContext(ctx, "UPDATE jobs SET last_error='stale',available_at=? WHERE completed_at IS NULL AND created_at<?", time.Now().UTC().Format(time.RFC3339Nano), before.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
