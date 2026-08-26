package auth

import (
	"context"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"testing"
	"time"
)

func testDB(t *testing.T) *storage.DB {
	t.Helper()
	db, e := storage.Open(context.Background(), "file::memory:?cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = storage.Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	return db
}
func TestPasswordRoundTrip(t *testing.T) {
	h := HashPassword("secret")
	if !CheckPassword(h, "secret") {
		t.Fatal("password mismatch")
	}
	if CheckPassword(h, "wrong") {
		t.Fatal("wrong password accepted")
	}
	if h == HashPassword("other") {
		t.Fatal("hash collision")
	}
}
func TestSessionLifecycle(t *testing.T) {
	db := testDB(t)
	u := repository.UserRepo{DB: db.SQL}
	s := repository.SessionRepo{DB: db.SQL}
	svc := Service{Users: u, Sessions: s, TTL: time.Hour}
	user, e := svc.Register(context.Background(), "officer", "pw", domain.RoleOfficer)
	if e != nil {
		t.Fatal(e)
	}
	ss, got, e := svc.Login(context.Background(), "officer", "pw")
	if e != nil || got.ID != user.ID {
		t.Fatalf("login %v", e)
	}
	if _, e = svc.Authenticate(context.Background(), ss.ID); e != nil {
		t.Fatal(e)
	}
	if e = svc.Logout(context.Background(), ss.ID); e != nil {
		t.Fatal(e)
	}
	if _, e = svc.Authenticate(context.Background(), ss.ID); e == nil {
		t.Fatal("revoked session accepted")
	}
}
func TestExpiredSession(t *testing.T) {
	db := testDB(t)
	u := repository.UserRepo{DB: db.SQL}
	s := repository.SessionRepo{DB: db.SQL}
	svc := Service{Users: u, Sessions: s, TTL: -time.Second}
	_, e := svc.Register(context.Background(), "review", "pw", domain.RoleReviewer)
	if e != nil {
		t.Fatal(e)
	}
	ss, _, e := svc.Login(context.Background(), "review", "pw")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = svc.Authenticate(context.Background(), ss.ID); e == nil {
		t.Fatal("expired accepted")
	}
}
func TestDuplicateUser(t *testing.T) {
	db := testDB(t)
	u := repository.UserRepo{DB: db.SQL}
	svc := Service{Users: u, TTL: time.Hour}
	if _, e := svc.Register(context.Background(), "same", "pw", domain.RoleOfficer); e != nil {
		t.Fatal(e)
	}
	if _, e := svc.Register(context.Background(), "same", "pw", domain.RoleOfficer); e == nil {
		t.Fatal("duplicate accepted")
	}
}
