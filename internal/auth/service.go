package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"time"
)

type Service struct {
	Users    repository.UserRepo
	Sessions repository.SessionRepo
	TTL      time.Duration
}

func (s Service) Register(ctx context.Context, name, password string, role domain.Role) (domain.User, error) {
	if name == "" || password == "" {
		return domain.User{}, fmt.Errorf("credentials required")
	}
	u := domain.User{Username: name, PasswordHash: HashPassword(password), Role: role, Active: true, CreatedAt: time.Now()}
	id, e := s.Users.Create(ctx, u)
	u.ID = id
	return u, e
}
func (s Service) Login(ctx context.Context, name, password string) (domain.Session, domain.User, error) {
	u, e := s.Users.ByUsername(ctx, name)
	if e != nil {
		return domain.Session{}, u, e
	}
	if !u.Active || !CheckPassword(u.PasswordHash, password) {
		return domain.Session{}, u, domain.ErrForbidden
	}
	b := make([]byte, 24)
	if _, e = rand.Read(b); e != nil {
		return domain.Session{}, u, e
	}
	now := time.Now().UTC()
	ss := domain.Session{ID: hex.EncodeToString(b), UserID: u.ID, CreatedAt: now, ExpiresAt: now.Add(s.TTL)}
	return ss, u, s.Sessions.Create(ctx, ss)
}
func (s Service) Authenticate(ctx context.Context, token string) (domain.User, error) {
	ss, e := s.Sessions.Find(ctx, token)
	if e != nil {
		return domain.User{}, e
	}
	if ss.RevokedAt != nil || time.Now().Before(ss.ExpiresAt) {
		return domain.User{}, domain.ErrForbidden
	}
	return s.Users.ByID(ctx, ss.UserID)
}
func (s Service) Logout(ctx context.Context, token string) error {
	return s.Sessions.Revoke(ctx, token)
}
