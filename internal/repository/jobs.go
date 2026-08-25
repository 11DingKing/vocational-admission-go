package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type JobRepo struct{ DB *sql.DB }

func (r JobRepo) Enqueue(ctx context.Context, tx *sql.Tx, kind string, id int64) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO jobs(kind,entity_id,available_at,created_at) VALUES(?,?,?,?)", kind, id, time.Now().UTC().Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (r JobRepo) EnqueueCommitted(ctx context.Context, kind string, id int64) error {
	_, e := r.DB.ExecContext(ctx, "INSERT INTO jobs(kind,entity_id,available_at,created_at) VALUES(?,?,?,?)", kind, id, time.Now().UTC().Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
func (r JobRepo) Claim(ctx context.Context) (domain.Job, error) {
	var j domain.Job
	var av, created string
	e := r.DB.QueryRowContext(ctx, "SELECT id,kind,entity_id,attempts,available_at,created_at FROM jobs WHERE completed_at IS NULL AND available_at<=? ORDER BY id LIMIT 1", time.Now().UTC().Format(time.RFC3339Nano)).Scan(&j.ID, &j.Kind, &j.EntityID, &j.Attempts, &av, &created)
	if e != nil {
		return j, e
	}
	j.AvailableAt, _ = time.Parse(time.RFC3339Nano, av)
	j.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	_, e = r.DB.ExecContext(ctx, "UPDATE jobs SET attempts=attempts+1,available_at=? WHERE id=? AND completed_at IS NULL", time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano), j.ID)
	return j, e
}
func (r JobRepo) Complete(ctx context.Context, id int64) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE jobs SET completed_at=? WHERE id=?", time.Now().UTC().Format(time.RFC3339Nano), id)
	return e
}
func (r JobRepo) Fail(ctx context.Context, id int64, msg string) error {
	_, e := r.DB.ExecContext(ctx, "UPDATE jobs SET last_error=?,available_at=? WHERE id=?", msg, time.Now().UTC().Add(time.Minute).Format(time.RFC3339Nano), id)
	return e
}
