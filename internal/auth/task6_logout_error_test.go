package auth

import (
	"context"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"testing"
	"time"
)

func TestLogoutPropagatesRevocationFailure(t *testing.T) {
	db, e := storage.Open(context.Background(), "file:logout?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	if e = storage.Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	if _, e = db.SQL.Exec("INSERT INTO users(username,password_hash,role,created_at) VALUES('u',?,'officer',datetime('now'))", HashPassword("pw")); e != nil {
		t.Fatal(e)
	}
	s := Service{Users: repository.UserRepo{DB: db.SQL}, Sessions: repository.SessionRepo{DB: db.SQL}, TTL: time.Hour}
	ss, _, e := s.Login(context.Background(), "u", "pw")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = db.SQL.Exec("CREATE TRIGGER fail_revoke BEFORE UPDATE OF revoked_at ON sessions BEGIN SELECT RAISE(ABORT,'revoke down'); END"); e != nil {
		t.Fatal(e)
	}
	if e = s.Logout(context.Background(), ss.ID); e == nil {
		t.Fatal("revocation failure hidden")
	}
}
