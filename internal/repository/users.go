package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type UserRepo struct{ DB *sql.DB }

func (r UserRepo) Create(ctx context.Context, u domain.User) (int64, error) {
	res, e := r.DB.ExecContext(ctx, "INSERT INTO users(username,password_hash,role,active,created_at) VALUES(?,?,?,?,?)", u.Username, u.PasswordHash, u.Role, boolInt(u.Active), u.CreatedAt.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, fmt.Errorf("create user: %w", e)
	}
	return res.LastInsertId()
}
func (r UserRepo) ByUsername(ctx context.Context, name string) (domain.User, error) {
	var u domain.User
	var active int
	var created string
	e := r.DB.QueryRowContext(ctx, "SELECT id,username,password_hash,role,active,created_at FROM users WHERE username=?", name).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &active, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return u, domain.ErrNotFound
	}
	if e != nil {
		return u, e
	}
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return u, nil
}
func (r UserRepo) ByID(ctx context.Context, id int64) (domain.User, error) {
	var u domain.User
	var active int
	var created string
	e := r.DB.QueryRowContext(ctx, "SELECT id,username,password_hash,role,active,created_at FROM users WHERE id=?", id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &active, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return u, domain.ErrNotFound
	}
	if e != nil {
		return u, e
	}
	u.Active = active == 1
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	return u, nil
}
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
