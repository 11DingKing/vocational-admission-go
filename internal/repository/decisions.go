package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"time"
)

func ListDecisions(ctx context.Context, db *sql.DB, applicationID int64) ([]domain.Decision, error) {
	rows, e := db.QueryContext(ctx, "SELECT id,application_id,actor_id,from_status,to_status,reason,request_id,created_at FROM decisions WHERE application_id=? ORDER BY id", applicationID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.Decision
	for rows.Next() {
		var d domain.Decision
		var ts string
		if e = rows.Scan(&d.ID, &d.ApplicationID, &d.ActorID, &d.FromStatus, &d.ToStatus, &d.Reason, &d.RequestID, &ts); e != nil {
			return nil, e
		}
		d.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, d)
	}
	return out, rows.Err()
}
func ListAudit(ctx context.Context, db *sql.DB, entity string, id int64) ([]domain.AuditEvent, error) {
	rows, e := db.QueryContext(ctx, "SELECT id,actor_id,entity,entity_id,action,outcome,request_id,details,created_at FROM audit_events WHERE entity=? AND entity_id=? ORDER BY id", entity, id)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []domain.AuditEvent
	for rows.Next() {
		var a domain.AuditEvent
		var ts string
		if e = rows.Scan(&a.ID, &a.ActorID, &a.Entity, &a.EntityID, &a.Action, &a.Outcome, &a.RequestID, &a.Details, &ts); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, a)
	}
	return out, rows.Err()
}
