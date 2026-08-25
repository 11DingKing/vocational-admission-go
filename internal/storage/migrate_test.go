package storage

import (
	"context"
	"database/sql"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"os"
	"testing"
	"time"
)

func openTest(t *testing.T) *DB {
	t.Helper()
	db, e := Open(context.Background(), "file:test-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
func TestMigrationCreatesTables(t *testing.T) {
	db := openTest(t)
	rows, e := db.SQL.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if e != nil {
		t.Fatal(e)
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var n string
		_ = rows.Scan(&n)
		seen[n] = true
	}
	for _, n := range []string{"users", "sessions", "province_rules", "plans", "major_groups", "applications", "decisions", "audit_events", "jobs"} {
		if !seen[n] {
			t.Errorf("missing %s", n)
		}
	}
}
func TestMigrationIdempotent(t *testing.T) {
	db := openTest(t)
	if e := Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	var n int
	if e := db.SQL.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		t.Fatal(e)
	}
	if n != 5 {
		t.Fatalf("versions=%d", n)
	}
}
func TestRestartRecovery(t *testing.T) {
	f, err := os.CreateTemp("", "admission-recovery-*.db")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	_ = f.Close()
	defer os.Remove(path)
	db, e := Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	if e = Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	u := repository.UserRepo{DB: db.SQL}
	if _, e = u.Create(context.Background(), domain.User{Username: "persist", PasswordHash: "h", Role: domain.RoleOfficer, Active: true, CreatedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	_ = db.Close()
	db2, e := Open(context.Background(), path)
	if e != nil {
		t.Fatal(e)
	}
	defer db2.Close()
	if e = Migrate(context.Background(), db2.SQL); e != nil {
		t.Fatal(e)
	}
	if _, e = (repository.UserRepo{DB: db2.SQL}).ByUsername(context.Background(), "persist"); e != nil {
		t.Fatal(e)
	}
}
func TestTxRollback(t *testing.T) {
	db := openTest(t)
	e := db.Tx(context.Background(), func(tx *sql.Tx) error {
		_, e := tx.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES('rollback','x','officer',datetime('now'))")
		if e != nil {
			return e
		}
		return sql.ErrTxDone
	})
	if e == nil {
		t.Fatal("expected rollback error")
	}
	var n int
	_ = db.SQL.QueryRow("SELECT COUNT(*) FROM users WHERE username='rollback'").Scan(&n)
	if n != 0 {
		t.Fatal("row committed")
	}
}
func TestSnapshot(t *testing.T) {
	db := openTest(t)
	v, e := Snapshot(context.Background(), db.SQL)
	if e != nil || v == "" {
		t.Fatalf("snapshot %q %v", v, e)
	}
}
