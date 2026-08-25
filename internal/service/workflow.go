package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
)

func (s AdmissionService) PublishPlan(ctx context.Context, id int64, actor int64, req string) error {
	p, e := s.Plans.ByID(ctx, id)
	if e != nil {
		return e
	}
	if ok := p.Status.CanTransition(domain.PlanPublished); !ok {
		return domain.ErrInvalidState
	}
	return s.Plans.Transition(ctx, id, p.Status, domain.PlanPublished, p.Version)
}
func (s AdmissionService) ClosePlan(ctx context.Context, id int64, actor int64, req string) error {
	return s.DB.Tx(ctx, func(tx *sql.Tx) error {
		p, e := s.Plans.ByID(ctx, id)
		if e != nil {
			return e
		}
		if !p.Status.CanTransition(domain.PlanClosed) {
			return domain.ErrInvalidState
		}
		// Block closing while any application has not reached a terminal
		// decision. Both freshly submitted and currently reviewing
		// applications must be processed first; once the plan is closed it
		// can no longer transition, so leaving such applications behind would
		// strand them permanently.
		pending, e := s.Apps.CountPending(ctx, tx, id)
		if e != nil {
			return e
		}
		if pending > 0 {
			return fmt.Errorf("%w: %d applications still pending", domain.ErrConflict, pending)
		}
		if e = s.Plans.Transition(ctx, id, p.Status, domain.PlanClosed, p.Version); e != nil {
			return e
		}
		return s.Audits.Record(ctx, tx, actor, "plan", id, "close", "success", req, map[string]any{"version": p.Version})
	})
}
func (s AdmissionService) Withdraw(ctx context.Context, id int64, actor int64, req string) error {
	return s.Decide(ctx, id, domain.ApplicationWithdrawn, "withdrawn by officer", actor, req)
}
func (s AdmissionService) Admit(ctx context.Context, id int64, actor int64, req string) error {
	return s.Decide(ctx, id, domain.ApplicationAdmitted, "qualified by review", actor, req)
}
func (s AdmissionService) Reject(ctx context.Context, id int64, reason string, actor int64, req string) error {
	if reason == "" {
		reason = "not qualified"
	}
	return s.Decide(ctx, id, domain.ApplicationRejected, reason, actor, req)
}
