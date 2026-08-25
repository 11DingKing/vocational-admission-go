package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type RuleRepo struct{ DB *sql.DB }

func (r RuleRepo) Upsert(ctx context.Context, p domain.ProvinceRule) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO province_rules(province,batch,min_score,rank_limit,allow_transfer,version,updated_at) VALUES(?,?,?,?,?,?,?) ON CONFLICT(province,batch) DO UPDATE SET min_score=excluded.min_score,rank_limit=excluded.rank_limit,allow_transfer=excluded.allow_transfer,version=province_rules.version+1,updated_at=excluded.updated_at", p.Province, p.Batch, p.MinScore, p.RankLimit, boolInt(p.AllowTransfer), p.Version, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (r RuleRepo) Get(ctx context.Context, province string, batch domain.Batch) (domain.ProvinceRule, error) {
	var p domain.ProvinceRule
	var transfer int
	var updated string
	e := r.DB.QueryRowContext(ctx, "SELECT id,province,batch,min_score,rank_limit,allow_transfer,version,updated_at FROM province_rules WHERE province=? AND batch=?", province, batch).Scan(&p.ID, &p.Province, &p.Batch, &p.MinScore, &p.RankLimit, &transfer, &p.Version, &updated)
	p.AllowTransfer = transfer == 1
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, e
}
func (r RuleRepo) List(ctx context.Context) ([]domain.ProvinceRule, error) {
	rows, e := r.DB.QueryContext(ctx, "SELECT id,province,batch,min_score,rank_limit,allow_transfer,version,updated_at FROM province_rules ORDER BY province,batch")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.ProvinceRule
	for rows.Next() {
		var p domain.ProvinceRule
		var t int
		var ts string
		if e = rows.Scan(&p.ID, &p.Province, &p.Batch, &p.MinScore, &p.RankLimit, &t, &p.Version, &ts); e != nil {
			return nil, e
		}
		p.AllowTransfer = t == 1
		p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, p)
	}
	return out, rows.Err()
}
