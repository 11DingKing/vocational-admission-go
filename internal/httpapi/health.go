package httpapi

import (
	"net/http"
	"time"
)

type Health struct {
	Started time.Time
	Version string
}

func (h Health) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	age := time.Since(h.Started)
	write(w, http.StatusOK, map[string]any{"status": "ok", "version": h.Version, "uptime_seconds": int(age.Seconds())})
}
func Method(next http.Handler, allowed ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range allowed {
			if r.Method == m {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("Allow", allowed[0])
		write(w, http.StatusMethodNotAllowed, nil)
	})
}
