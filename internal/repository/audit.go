package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

type AuditRepo struct{ DB *sql.DB }

func (r AuditRepo) Record(ctx context.Context, tx *sql.Tx, e domain.AuditEvent) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO audit_events(actor_id,entity,entity_id,action,outcome,request_id,details,created_at) VALUES(?,?,?,?,?,?,?,?)", e.ActorID, e.Entity, e.EntityID, e.Action, e.Outcome, e.RequestID, e.Details, e.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}
func (r AuditRepo) Decide(ctx context.Context, tx *sql.Tx, d domain.Decision) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO decisions(application_id,actor_id,from_status,to_status,reason,request_id,created_at) VALUES(?,?,?,?,?,?,?)", d.ApplicationID, d.ActorID, d.FromStatus, d.ToStatus, d.Reason, d.RequestID, d.CreatedAt.UTC().Format(time.RFC3339Nano))
	return err
}
