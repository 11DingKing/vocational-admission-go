package service

import (
	"context"
	"github.com/11DingKing/vocational-admission-go/internal/audit"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"testing"
	"time"
)

func TestSubmittedApplicationKeepsPlanOpen(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, "file:task1-close?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = storage.Migrate(ctx, db.SQL); err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{"('officer','x','officer',datetime('now'))", "('reviewer','x','reviewer',datetime('now'))"} {
		if _, err = db.SQL.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES" + row); err != nil {
			t.Fatal(err)
		}
	}
	plans := repository.PlanRepo{DB: db.SQL}
	pid, err := plans.Create(ctx, domain.AdmissionPlan{Year: 2026, Province: "浙江", Name: "本科计划", Status: domain.PlanPublished, TotalCapacity: 2, CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	gid, err := plans.AddGroup(ctx, domain.MajorGroup{PlanID: pid, Code: "G1", Name: "智能制造", Capacity: 2})
	if err != nil {
		t.Fatal(err)
	}
	svc := AdmissionService{DB: db, Plans: plans, Apps: repository.ApplicationRepo{DB: db.SQL}, Audits: audit.Logger{Repo: repository.AuditRepo{DB: db.SQL}}, Jobs: repository.JobRepo{DB: db.SQL}}
	if _, err = svc.Submit(ctx, domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260101", Score: 620, Rank: 10, IdempotencyKey: "task1-key"}, 1, "task1-submit"); err != nil {
		t.Fatal(err)
	}
	if err = svc.ClosePlan(ctx, pid, 1, "task1-close"); err == nil {
		t.Fatal("plan closed while a submitted application was pending")
	}
	plan, err := plans.ByID(ctx, pid)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != domain.PlanPublished {
		t.Fatalf("plan status changed to %s", plan.Status)
	}
}
