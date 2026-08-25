package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"testing"
	"time"
)

func repoDB(t *testing.T) *storage.DB {
	db, e := storage.Open(context.Background(), "file:repo-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = storage.Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
func TestPlanAndGroup(t *testing.T) {
	db := repoDB(t)
	r := PlanRepo{DB: db.SQL}
	id, e := r.Create(context.Background(), domain.AdmissionPlan{Year: 2026, Province: "浙江", Name: "计划", Status: domain.PlanDraft, TotalCapacity: 4, CreatedAt: time.Now()})
	if e != nil {
		t.Fatal(e)
	}
	gid, e := r.AddGroup(context.Background(), domain.MajorGroup{PlanID: id, Code: "A", Name: "智能制造", Capacity: 2})
	if e != nil {
		t.Fatal(e)
	}
	g, e := r.Group(context.Background(), gid)
	if e != nil || g.Capacity != 2 {
		t.Fatalf("group %#v %v", g, e)
	}
	p, e := r.ByID(context.Background(), id)
	if e != nil {
		t.Fatal(e)
	}
	if e = r.Transition(context.Background(), id, p.Status, domain.PlanLocked, p.Version); e != nil {
		t.Fatal(e)
	}
}
func TestReserveCapacity(t *testing.T) {
	db := repoDB(t)
	r := PlanRepo{DB: db.SQL}
	id, _ := r.Create(context.Background(), domain.AdmissionPlan{Year: 2026, Province: "江苏", Name: "计划", Status: domain.PlanPublished, TotalCapacity: 1, CreatedAt: time.Now()})
	gid, _ := r.AddGroup(context.Background(), domain.MajorGroup{PlanID: id, Code: "A", Name: "专业", Capacity: 1})
	if e := db.Tx(context.Background(), func(tx *sql.Tx) error { return r.Reserve(context.Background(), tx, id, gid) }); e != nil {
		t.Fatal(e)
	}
	if e := db.Tx(context.Background(), func(tx *sql.Tx) error { return r.Reserve(context.Background(), tx, id, gid) }); e == nil {
		t.Fatal("capacity exceeded not detected")
	}
}
func TestRules(t *testing.T) {
	db := repoDB(t)
	r := RuleRepo{DB: db.SQL}
	p := domain.ProvinceRule{Province: "北京", Batch: domain.BatchRegular, MinScore: 450, RankLimit: 1000, AllowTransfer: true, Version: 1}
	if e := r.Upsert(context.Background(), p); e != nil {
		t.Fatal(e)
	}
	got, e := r.Get(context.Background(), "北京", domain.BatchRegular)
	if e != nil || got.MinScore != 450 || !got.AllowTransfer {
		t.Fatalf("rule %#v %v", got, e)
	}
	p.MinScore = 500
	if e = r.Upsert(context.Background(), p); e != nil {
		t.Fatal(e)
	}
	got, e = r.Get(context.Background(), "北京", domain.BatchRegular)
	if e != nil || got.MinScore != 500 {
		t.Fatal(e)
	}
}
