package service

import (
	"context"
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
	p, e := s.Plans.ByID(ctx, id)
	if e != nil {
		return e
	}
	if p.Status != domain.PlanPublished {
		return domain.ErrInvalidState
	}
	sum, e := s.List(ctx, id, string(domain.ApplicationSubmitted), 1, 0)
	if e != nil {
		return e
	}
	if len(sum) > 0 {
		return fmt.Errorf("%w: pending applications", domain.ErrConflict)
	}
	return s.Plans.Transition(ctx, id, p.Status, domain.PlanClosed, p.Version)
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
