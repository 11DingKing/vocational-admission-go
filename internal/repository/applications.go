package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type ApplicationRepo struct{ DB *sql.DB }

func (r ApplicationRepo) Create(ctx context.Context, tx *sql.Tx, a domain.Application) (int64, error) {
	b, _ := json.Marshal(a.Preferences)
	res, e := tx.ExecContext(ctx, "INSERT INTO applications(plan_id,major_group_id,student_no,score,rank,preferences,status,transfer_accepted,idempotency_key,submitted_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)", a.PlanID, a.MajorGroupID, a.StudentNo, a.Score, a.Rank, string(b), a.Status, boolInt(a.TransferAccepted), a.IdempotencyKey, a.SubmittedAt.UTC().Format(time.RFC3339Nano), a.UpdatedAt.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, fmt.Errorf("create application: %w", e)
	}
	id, e := res.LastInsertId()
	return id, e
}
func (r ApplicationRepo) ByID(ctx context.Context, id int64) (domain.Application, error) {
	var a domain.Application
	var pref, status, submitted, updated string
	var transfer int
	e := r.DB.QueryRowContext(ctx, "SELECT id,plan_id,major_group_id,student_no,score,rank,preferences,status,transfer_accepted,idempotency_key,version,submitted_at,updated_at FROM applications WHERE id=?", id).Scan(&a.ID, &a.PlanID, &a.MajorGroupID, &a.StudentNo, &a.Score, &a.Rank, &pref, &status, &transfer, &a.IdempotencyKey, &a.Version, &submitted, &updated)
	if errors.Is(e, sql.ErrNoRows) {
		return a, domain.ErrNotFound
	}
	if e != nil {
		return a, e
	}
	a.Status = domain.ApplicationStatus(status)
	a.TransferAccepted = transfer == 1
	_ = json.Unmarshal([]byte(pref), &a.Preferences)
	a.SubmittedAt, _ = time.Parse(time.RFC3339Nano, submitted)
	a.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return a, nil
}
func (r ApplicationRepo) ByKey(ctx context.Context, key string) (domain.Application, error) {
	var id int64
	e := r.DB.QueryRowContext(ctx, "SELECT id FROM applications WHERE idempotency_key=?", key).Scan(&id)
	if errors.Is(e, sql.ErrNoRows) {
		return domain.Application{}, domain.ErrNotFound
	}
	if e != nil {
		return domain.Application{}, e
	}
	return r.ByID(ctx, id)
}
func (r ApplicationRepo) Transition(ctx context.Context, id int64, from, to domain.ApplicationStatus, version int) error {
	res, e := r.DB.ExecContext(ctx, "UPDATE applications SET status=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?", to, time.Now().UTC().Format(time.RFC3339Nano), id, from, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r ApplicationRepo) TransitionCommitted(ctx context.Context, id int64, from, to domain.ApplicationStatus, version int) error {
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
func (r ApplicationRepo) List(ctx context.Context, planID int64, status string, limit, offset int) ([]domain.Application, error) {
	q := "SELECT id FROM applications WHERE plan_id=?"
	args := []any{planID}
	if status != "" {
		q += " AND status=?"
		args = append(args, status)
	}
	q += " ORDER BY score DESC,rank ASC,id ASC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, e := r.DB.QueryContext(ctx, q, args...)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Application
	for rows.Next() {
		var id int64
		if e = rows.Scan(&id); e != nil {
			return nil, e
		}
		a, e := r.ByID(ctx, id)
		if e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
