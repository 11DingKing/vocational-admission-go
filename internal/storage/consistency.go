package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type ConsistencyIssue struct {
	Entity  string
	ID      int64
	Message string
}

func CheckConsistency(ctx context.Context, db *sql.DB) ([]ConsistencyIssue, error) {
	rows, e := db.QueryContext(ctx, "SELECT id,plan_id,major_group_id FROM applications")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []ConsistencyIssue
	for rows.Next() {
		var id, p, g int64
		if e = rows.Scan(&id, &p, &g); e != nil {
			return nil, e
		}
		var ok int
		if e = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM major_groups WHERE id=? AND plan_id=?", g, p).Scan(&ok); e != nil {
			return nil, e
		}
		if ok != 1 {
			out = append(out, ConsistencyIssue{Entity: "application", ID: id, Message: fmt.Sprintf("group %d not in plan %d", g, p)})
		}
	}
	return out, rows.Err()
}
func Vacuum(ctx context.Context, db *sql.DB) error { _, e := db.ExecContext(ctx, "VACUUM"); return e }
