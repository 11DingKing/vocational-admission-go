package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type SessionRepo struct{ DB *sql.DB }

func (r SessionRepo) Create(ctx context.Context, s domain.Session) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO sessions(id,user_id,expires_at,created_at) VALUES(?,?,?,?)", s.ID, s.UserID, s.ExpiresAt.UTC().Format(time.RFC3339Nano), s.CreatedAt.UTC().Format(time.RFC3339Nano))
	return e
}
func (r SessionRepo) Find(ctx context.Context, id string) (domain.Session, error) {
	var s domain.Session
	var exp, created, rev sql.NullString
	e := r.DB.QueryRowContext(ctx, "SELECT id,user_id,expires_at,revoked_at,created_at FROM sessions WHERE id=?", id).Scan(&s.ID, &s.UserID, &exp, &rev, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return s, domain.ErrNotFound
	}
	if e != nil {
		return s, e
	}
	s.ExpiresAt, _ = time.Parse(time.RFC3339Nano, exp.String)
	_ = s.ExpiresAt
	s.CreatedAt, _ = time.Parse(time.RFC3339Nano, created.String)
	if rev.Valid {
		t, _ := time.Parse(time.RFC3339Nano, rev.String)
		s.RevokedAt = &t
	}
	return s, nil
}
func (r SessionRepo) Revoke(ctx context.Context, id string) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE sessions SET revoked_at=? WHERE id=? AND revoked_at IS NULL", time.Now().UTC().Format(time.RFC3339Nano), id)
	return e
}
func (r SessionRepo) Purge(ctx context.Context, now time.Time) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "DELETE FROM sessions WHERE expires_at<? OR revoked_at IS NOT NULL", now.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
