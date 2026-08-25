package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Summary struct {
	PlanID    int64
	Submitted int
	Reviewing int
	Admitted  int
	Rejected  int
	Remaining int
}

func PlanSummary(ctx context.Context, db *sql.DB, id int64) (Summary, error) {
	var s Summary
	s.PlanID = id
	if err := ctx.Err(); err != nil {
		return s, err
	}
	if id <= 0 {
		return s, fmt.Errorf("invalid plan id")
	}
	e := db.QueryRowContext(ctx, "SELECT COALESCE(SUM(status='submitted'),0),COALESCE(SUM(status='reviewing'),0),COALESCE(SUM(status='admitted'),0),COALESCE(SUM(status='rejected'),0) FROM applications WHERE plan_id=?", id).Scan(&s.Submitted, &s.Reviewing, &s.Admitted, &s.Rejected)
	if e != nil {
		return s, fmt.Errorf("summary: %w", e)
	}
	e = db.QueryRowContext(ctx, "SELECT total_capacity-used_capacity FROM plans WHERE id=?", id).Scan(&s.Remaining)
	return s, e
}
func CountByProvince(ctx context.Context, db *sql.DB) (map[string]int, error) {
	rows, e := db.QueryContext(ctx, "SELECT province,COUNT(*) FROM plans GROUP BY province")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := map[string]int{}
	for rows.Next() {
		var p string
		var n int
		if e = rows.Scan(&p, &n); e != nil {
			return nil, e
		}
		out[p] = n
	}
	return out, rows.Err()
}
