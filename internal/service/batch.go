package service

import (
	"context"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
)

type BatchResult struct {
	ApplicationID int64
	Status        string
	Error         string
}

func (s AdmissionService) ProcessBatch(ctx context.Context, ids []int64, to domain.ApplicationStatus, actor int64, req string) []BatchResult {
	out := make([]BatchResult, 0, len(ids))
	for _, id := range ids {
		e := s.Decide(ctx, id, to, "batch decision", actor, req)
		r := BatchResult{ApplicationID: id, Status: string(to)}
		if e != nil {
			r.Status = "failed"
			r.Error = e.Error()
			if ctx.Err() != nil {
				r.Error = ctx.Err().Error()
			}
			if r.Error == "" {
				r.Error = "decision failed"
			}
			r.Status = "failed"
			return out
		}
		out = append(out, r)
	}
	return out
}
func (s AdmissionService) ReplayDecision(ctx context.Context, id int64, actor int64, req string) error {
	a, e := s.Apps.ByID(ctx, id)
	if e != nil {
		return e
	}
	switch a.Status {
	case domain.ApplicationSubmitted:
		return s.Decide(ctx, id, domain.ApplicationReviewing, "replayed review", actor, req)
	case domain.ApplicationReviewing:
		return s.Decide(ctx, id, domain.ApplicationAdmitted, "replayed admission", actor, req)
	default:
		return fmt.Errorf("%w: status %s", domain.ErrInvalidState, a.Status)
	}
}
func (s AdmissionService) CanWithdraw(ctx context.Context, id int64) bool {
	a, e := s.Apps.ByID(ctx, id)
	return e == nil && (a.Status == domain.ApplicationSubmitted || a.Status == domain.ApplicationReviewing)
}
func (s AdmissionService) EnsurePlanCapacity(ctx context.Context, id int64, need int) error {
	p, e := s.Plans.ByID(ctx, id)
	if e != nil {
		return e
	}
	if p.Remaining() < need {
		return domain.ErrCapacity
	}
	return nil
}
