package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func Snapshot(ctx context.Context, db *sql.DB) (string, error) {
	var now string
	if e := db.QueryRowContext(ctx, "SELECT datetime('now')").Scan(&now); e != nil {
		return "", e
	}
	return fmt.Sprintf("snapshot-%s", now), nil
}
func PurgeAudit(ctx context.Context, db *sql.DB, before time.Time) (int64, error) {
	res, e := db.ExecContext(ctx, "DELETE FROM audit_events WHERE created_at<?", before.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
