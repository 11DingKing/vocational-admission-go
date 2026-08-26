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

func TestNotificationFailureRollsBackAdmission(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, "file:task2-submit?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = storage.Migrate(ctx, db.SQL); err != nil {
		t.Fatal(err)
	}
	if _, err = db.SQL.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES('officer','x','officer',datetime('now'))"); err != nil {
		t.Fatal(err)
	}
	if _, err = db.SQL.Exec("CREATE TRIGGER fail_jobs BEFORE INSERT ON jobs BEGIN SELECT RAISE(ABORT,'jobs down'); END"); err != nil {
		t.Fatal(err)
	}
	plans := repository.PlanRepo{DB: db.SQL}
	pid, err := plans.Create(ctx, domain.AdmissionPlan{Year: 2026, Province: "江苏", Name: "普通批", Status: domain.PlanPublished, TotalCapacity: 1, CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	gid, err := plans.AddGroup(ctx, domain.MajorGroup{PlanID: pid, Code: "G1", Name: "软件工程", Capacity: 1})
	if err != nil {
		t.Fatal(err)
	}
	svc := AdmissionService{DB: db, Plans: plans, Apps: repository.ApplicationRepo{DB: db.SQL}, Audits: audit.Logger{Repo: repository.AuditRepo{DB: db.SQL}}, Jobs: repository.JobRepo{DB: db.SQL}}
	if _, err = svc.Submit(ctx, domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260201", Score: 630, Rank: 8, IdempotencyKey: "task2-key"}, 1, "task2-submit"); err == nil {
		t.Fatal("notification failure was hidden")
	}
	var apps, used int
	if err = db.SQL.QueryRow("SELECT COUNT(*) FROM applications WHERE idempotency_key='task2-key'").Scan(&apps); err != nil {
		t.Fatal(err)
	}
	if err = db.SQL.QueryRow("SELECT used_capacity FROM plans WHERE id=?", pid).Scan(&used); err != nil {
		t.Fatal(err)
	}
	if apps != 0 || used != 0 {
		t.Fatalf("failed submission leaked application=%d capacity=%d", apps, used)
	}
}
