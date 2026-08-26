package service

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/audit"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"testing"
	"time"
)

func svc(t *testing.T) (AdmissionService, *storage.DB, int64, int64) {
	db, e := storage.Open(context.Background(), "file:svc-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = storage.Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	if _, e = db.SQL.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES('actor','h','officer',datetime('now'))"); e != nil {
		t.Fatal(e)
	}
	if _, e = db.SQL.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES('reviewer','h','reviewer',datetime('now'))"); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	r := repository.PlanRepo{DB: db.SQL}
	pid, e := r.Create(context.Background(), domain.AdmissionPlan{Year: 2026, Province: "上海", Name: "本科", Status: domain.PlanPublished, TotalCapacity: 2, CreatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	gid, e := r.AddGroup(context.Background(), domain.MajorGroup{PlanID: pid, Code: "G1", Name: "软件", Capacity: 1})
	if e != nil {
		t.Fatal(e)
	}
	s := AdmissionService{DB: db, Plans: r, Apps: repository.ApplicationRepo{DB: db.SQL}, Audits: audit.Logger{Repo: repository.AuditRepo{DB: db.SQL}}, Jobs: repository.JobRepo{DB: db.SQL}}
	return s, db, pid, gid
}
func TestSubmitIdempotent(t *testing.T) {
	s, _, pid, gid := svc(t)
	a := domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260001", Score: 600, Rank: 1, IdempotencyKey: "key-1"}
	got, e := s.Submit(context.Background(), a, 1, "r1")
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.Submit(context.Background(), a, 1, "r2")
	if e != nil || again.ID != got.ID {
		t.Fatalf("idempotency %v %#v", e, again)
	}
}
func TestSubmitCapacity(t *testing.T) {
	s, _, pid, gid := svc(t)
	for i := 0; i < 2; i++ {
		a := domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260" + string(rune('0'+i)) + "1", Score: 600, Rank: i + 1, IdempotencyKey: "cap-" + string(rune('a'+i))}
		e := error(nil)
		_, e = s.Submit(context.Background(), a, 1, "r")
		if i == 0 && e != nil {
			t.Fatal(e)
		}
		if i == 1 && e == nil {
			t.Fatal("second reservation accepted")
		}
	}
}
func TestDecisionFlow(t *testing.T) {
	s, _, pid, gid := svc(t)
	a, e := s.Submit(context.Background(), domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260002", Score: 700, Rank: 2, IdempotencyKey: "flow"}, 1, "r")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Decide(context.Background(), a.ID, domain.ApplicationAdmitted, "qualified", 2, "r2"); e == nil {
		t.Fatal("submitted admitted without review")
	}
	if e = s.Decide(context.Background(), a.ID, domain.ApplicationReviewing, "review", 2, "r3"); e != nil {
		t.Fatal(e)
	}
	if e = s.Admit(context.Background(), a.ID, 2, "r4"); e != nil {
		t.Fatal(e)
	}
}
func TestDecisionConflict(t *testing.T) {
	s, _, pid, gid := svc(t)
	a, e := s.Submit(context.Background(), domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260003", Score: 700, Rank: 3, IdempotencyKey: "conflict"}, 1, "r")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.Reject(context.Background(), a.ID, "low", 2, "r"); e == nil {
		t.Fatal("invalid transition accepted")
	}
}
func TestContextCancellation(t *testing.T) {
	s, _, pid, gid := svc(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e := s.Submit(ctx, domain.Application{PlanID: pid, MajorGroupID: gid, StudentNo: "20260004", Score: 700, Rank: 4, IdempotencyKey: "cancel"}, 1, "r")
	if e == nil {
		t.Fatal("cancel ignored")
	}
}

var _ *sql.Tx
