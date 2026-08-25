package middleware

import (
	"github.com/11DingKing/vocational-admission-go/internal/auth"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"log/slog"
	"net/http"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = "req-" + r.RemoteAddr
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				slog.Error("panic recovered", "error", v)
				http.Error(w, "internal error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
func Authenticate(a auth.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			if u, e := a.Authenticate(r.Context(), token[7:]); e == nil {
				next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), u)))
				return
			}
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})
}

var _ domain.Role
