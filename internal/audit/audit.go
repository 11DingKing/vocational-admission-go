package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"time"
)

type Logger struct{ Repo repository.AuditRepo }

func (l Logger) Record(ctx context.Context, tx *sql.Tx, actor int64, entity string, id int64, action, outcome, request string, details any) error {
	b, _ := json.Marshal(details)
	return l.Repo.Record(ctx, tx, domain.AuditEvent{ActorID: actor, Entity: entity, EntityID: id, Action: action, Outcome: outcome, RequestID: request, Details: string(b), CreatedAt: time.Now()})
}
