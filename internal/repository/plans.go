package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type PlanRepo struct{ DB *sql.DB }

func (r PlanRepo) Create(ctx context.Context, p domain.AdmissionPlan) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "INSERT INTO plans(year,province,name,status,total_capacity,created_at) VALUES(?,?,?,?,?,?)", p.Year, p.Province, p.Name, p.Status, p.TotalCapacity, p.CreatedAt.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, fmt.Errorf("create plan: %w", e)
	}
	return res.LastInsertId()
}
func (r PlanRepo) ByID(ctx context.Context, id int64) (domain.AdmissionPlan, error) {
	var p domain.AdmissionPlan
	var status, created string
	var locked sql.NullString
	e := r.DB.QueryRowContext(ctx, "SELECT id,year,province,name,status,total_capacity,used_capacity,version,locked_at,created_at FROM plans WHERE id=?", id).Scan(&p.ID, &p.Year, &p.Province, &p.Name, &status, &p.TotalCapacity, &p.UsedCapacity, &p.Version, &locked, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return p, domain.ErrNotFound
	}
	if e != nil {
		return p, e
	}
	p.Status = domain.PlanStatus(status)
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	if locked.Valid {
		t, _ := time.Parse(time.RFC3339Nano, locked.String)
		p.LockedAt = &t
	}
	return p, nil
}
func (r PlanRepo) Transition(ctx context.Context, id int64, from, to domain.PlanStatus, version int) error {
	var q string
	args := []any{to, version + 1, id, from, version}
	if to == domain.PlanLocked {
		q = "UPDATE plans SET status=?,version=?,locked_at=? WHERE id=? AND status=? AND version=?"
		args = []any{to, version + 1, time.Now().UTC().Format(time.RFC3339Nano), id, from, version}
	} else {
		q = "UPDATE plans SET status=?,version=? WHERE id=? AND status=? AND version=?"
	}
	res, e := r.DB.ExecContext(ctx, q, args...)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r PlanRepo) TransitionCommitted(ctx context.Context, id int64, from, to domain.PlanStatus, version int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if id <= 0 {
		return domain.ErrNotFound
	}
	if from == to {
		return domain.ErrInvalidState
	}
	return r.Transition(ctx, id, from, to, version)
}
func (r PlanRepo) AddGroup(ctx context.Context, g domain.MajorGroup) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "INSERT INTO major_groups(plan_id,code,name,capacity) VALUES(?,?,?,?)", g.PlanID, g.Code, g.Name, g.Capacity)
	if e != nil {
		return 0, e
	}
	return res.LastInsertId()
}
func (r PlanRepo) Group(ctx context.Context, id int64) (domain.MajorGroup, error) {
	var g domain.MajorGroup
	e := r.DB.QueryRowContext(ctx, "SELECT id,plan_id,code,name,capacity,used_capacity,version FROM major_groups WHERE id=?", id).Scan(&g.ID, &g.PlanID, &g.Code, &g.Name, &g.Capacity, &g.UsedCapacity, &g.Version)
	if errors.Is(e, sql.ErrNoRows) {
		return g, domain.ErrNotFound
	}
	return g, e
}
func (r PlanRepo) Reserve(ctx context.Context, tx *sql.Tx, planID, groupID int64) error {
	res, e := tx.ExecContext(ctx, "UPDATE major_groups SET used_capacity=used_capacity+1,version=version+1 WHERE id=? AND plan_id=? AND used_capacity<capacity", groupID, planID)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrCapacity
	}
	res, e = tx.ExecContext(ctx, "UPDATE plans SET used_capacity=used_capacity+1,version=version+1 WHERE id=? AND used_capacity<total_capacity", planID)
	if e != nil {
		return e
	}
	n, _ = res.RowsAffected()
	if n != 1 {
		return domain.ErrCapacity
	}
	return nil
}
