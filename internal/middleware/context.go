package middleware

import (
	"context"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"net/http"
)

type key string

const userKey key = "user"

func WithUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}
func UserFrom(ctx context.Context) (domain.User, bool) {
	u, ok := ctx.Value(userKey).(domain.User)
	return u, ok
}
func Require(next http.Handler, roles ...domain.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFrom(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		for _, role := range roles {
			if u.Role == role || u.Role == domain.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "forbidden", http.StatusForbidden)
	})
}
