package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
)

type ReportService struct{ DB *sql.DB }

func (r ReportService) PlanSummary(ctx context.Context, id int64) (repository.Summary, error) {
	return repository.PlanSummary(ctx, r.DB, id)
}
func (r ReportService) ProvinceCounts(ctx context.Context) (map[string]int, error) {
	return repository.CountByProvince(ctx, r.DB)
}
func (r ReportService) RecentAudit(ctx context.Context, limit int) ([]string, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT action FROM audit_events ORDER BY id DESC LIMIT ?", limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if e = rows.Scan(&v); e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
