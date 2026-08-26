package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/audit"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"time"
)

type AdmissionService struct {
	DB     *storage.DB
	Plans  repository.PlanRepo
	Apps   repository.ApplicationRepo
	Audits audit.Logger
	Jobs   repository.JobRepo
}

func (s AdmissionService) CreatePlan(ctx context.Context, p domain.AdmissionPlan, actor int64, req string) (int64, error) {
	if p.TotalCapacity <= 0 {
		return 0, fmt.Errorf("capacity must be positive")
	}
	p.Status = domain.PlanDraft
	p.CreatedAt = time.Now()
	return s.Plans.Create(ctx, p)
}
func (s AdmissionService) LockPlan(ctx context.Context, id int64, actor int64, req string) error {
	return s.DB.Tx(ctx, func(tx *sql.Tx) error {
		p, e := s.Plans.ByID(ctx, id)
		if e != nil {
			return e
		}
		if !p.Status.CanTransition(domain.PlanLocked) {
			return domain.ErrInvalidState
		}
		if e = s.Plans.Transition(ctx, id, p.Status, domain.PlanLocked, p.Version); e != nil {
			return e
		}
		return s.Audits.Record(ctx, tx, actor, "plan", id, "lock", "success", req, map[string]any{"version": p.Version})
	})
}
func (s AdmissionService) Submit(ctx context.Context, a domain.Application, actor int64, req string) (domain.Application, error) {
	if a.IdempotencyKey == "" {
		return a, domain.ErrDuplicate
	}
	if old, e := s.Apps.ByKey(ctx, a.IdempotencyKey); e == nil {
		return old, nil
	}
	plan, e := s.Plans.ByID(ctx, a.PlanID)
	if e != nil {
		return a, e
	}
	if plan.Status != domain.PlanPublished && plan.Status != domain.PlanLocked {
		return a, domain.ErrInvalidState
	}
	group, e := s.Plans.Group(ctx, a.MajorGroupID)
	if e != nil {
		return a, e
	}
	if group.PlanID != plan.ID {
		return a, domain.ErrConflict
	}
	a.Status = domain.ApplicationSubmitted
	a.SubmittedAt = time.Now()
	a.UpdatedAt = a.SubmittedAt
	var id int64
	e = s.DB.Tx(ctx, func(tx *sql.Tx) error {
		if e := s.Plans.Reserve(ctx, tx, a.PlanID, a.MajorGroupID); e != nil {
			return e
		}
		id, e = s.Apps.Create(ctx, tx, a)
		if e != nil {
			return e
		}
		if e = s.Audits.Record(ctx, tx, actor, "application", id, "submit", "success", req, a.StudentNo); e != nil {
			return e
		}
		return s.Jobs.Enqueue(ctx, tx, "notify_submission", id)
	})
	if e != nil {
		return a, e
	}
	return s.Apps.ByID(ctx, id)
}
func (s AdmissionService) Decide(ctx context.Context, id int64, to domain.ApplicationStatus, reason string, actor int64, req string) error {
	return s.DB.Tx(ctx, func(tx *sql.Tx) error {
		a, e := s.Apps.ByID(ctx, id)
		if e != nil {
			return e
		}
		if !a.Status.CanTransition(to) {
			return domain.ErrInvalidState
		}
		if to == domain.ApplicationAdmitted {
			rule := domain.ProvinceRule{MinScore: 0}
			if e = domain.ValidateScore(rule, a.Score, a.Rank); e != nil {
				return e
			}
		}
		if e = s.Apps.Transition(ctx, id, a.Status, to, a.Version); e != nil {
			return e
		}
		if e = s.Audits.Repo.Decide(ctx, tx, domain.Decision{ApplicationID: id, ActorID: actor, FromStatus: string(a.Status), ToStatus: string(to), Reason: reason, RequestID: req, CreatedAt: time.Now()}); e != nil {
			return e
		}
		return s.Audits.Record(ctx, tx, actor, "application", id, "decision", "success", req, reason)
	})
}
func (s AdmissionService) List(ctx context.Context, plan int64, status string, limit, offset int) ([]domain.Application, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.Apps.List(ctx, plan, status, limit, offset)
}

var _ = errors.Is
