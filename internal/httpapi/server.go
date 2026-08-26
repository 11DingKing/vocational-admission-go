package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/11DingKing/vocational-admission-go/internal/auth"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/middleware"
	"github.com/11DingKing/vocational-admission-go/internal/pagination"
	"github.com/11DingKing/vocational-admission-go/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	Auth       auth.Service
	Admissions service.AdmissionService
	Mux        *http.ServeMux
}

func New(a auth.Service, s service.AdmissionService) *Server {
	h := &Server{Auth: a, Admissions: s, Mux: http.NewServeMux()}
	h.routes()
	return h
}
func (h *Server) routes() {
	h.Mux.HandleFunc("/healthz", h.health)
	h.Mux.HandleFunc("/readyz", h.ready)
	h.Mux.HandleFunc("/v1/auth/register", h.register)
	h.Mux.HandleFunc("/v1/auth/login", h.login)
	h.Mux.Handle("/v1/auth/logout", middleware.Authenticate(h.Auth, http.HandlerFunc(h.logout)))
	h.Mux.Handle("/v1/plans", middleware.Authenticate(h.Auth, middleware.Require(http.HandlerFunc(h.plans), domain.RoleOfficer, domain.RoleReviewer)))
	h.Mux.Handle("/v1/applications", middleware.Authenticate(h.Auth, middleware.Require(http.HandlerFunc(h.apps), domain.RoleOfficer, domain.RoleReviewer)))
}
func (h *Server) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]string{"status": "ok"})
}
func (h *Server) ready(w http.ResponseWriter, r *http.Request) {
	if e := h.Admissions.DB.Ping(r.Context()); e != nil {
		write(w, 503, map[string]string{"status": "not_ready"})
		return
	}
	write(w, 200, map[string]string{"status": "ready"})
}
func (h *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in struct {
		Username, Password string
		Role               domain.Role
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		write(w, 400, nil)
		return
	}
	u, e := h.Auth.Register(r.Context(), in.Username, in.Password, in.Role)
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 201, u)
}
func (h *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		write(w, 405, nil)
		return
	}
	var in struct{ Username, Password string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		write(w, 400, nil)
		return
	}
	ss, u, e := h.Auth.Login(r.Context(), in.Username, in.Password)
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 200, map[string]any{"token": ss.ID, "expires_at": ss.ExpiresAt, "user": u})
}
func (h *Server) logout(w http.ResponseWriter, r *http.Request) {
	_ = h.Auth.Logout(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	write(w, 204, nil)
}
func (h *Server) plans(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var in struct {
			Year           int
			Province, Name string
			Capacity       int
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			write(w, 400, nil)
			return
		}
		u, _ := middleware.UserFrom(r.Context())
		id, e := h.Admissions.CreatePlan(r.Context(), domain.AdmissionPlan{Year: in.Year, Province: in.Province, Name: in.Name, TotalCapacity: in.Capacity}, u.ID, r.Header.Get("X-Request-ID"))
		if e != nil {
			writeErr(w, e)
			return
		}
		write(w, 201, map[string]any{"id": id})
		return
	}
	write(w, 405, nil)
}
func (h *Server) apps(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var in struct {
			PlanID, MajorGroupID int64
			StudentNo            string
			Score, Rank          int
			Preferences          []string
			TransferAccepted     bool
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			write(w, 400, nil)
			return
		}
		u, _ := middleware.UserFrom(r.Context())
		a, e := h.Admissions.Submit(r.Context(), domain.Application{PlanID: in.PlanID, MajorGroupID: in.MajorGroupID, StudentNo: in.StudentNo, Score: in.Score, Rank: in.Rank, Preferences: in.Preferences, TransferAccepted: in.TransferAccepted, IdempotencyKey: r.Header.Get("Idempotency-Key")}, u.ID, r.Header.Get("X-Request-ID"))
		if e != nil {
			writeErr(w, e)
			return
		}
		write(w, 201, a)
		return
	}
	if r.Method == "GET" {
		pid, _ := strconv.ParseInt(r.URL.Query().Get("plan_id"), 10, 64)
		q := pagination.Parse(atoi(r.URL.Query().Get("limit")), atoi(r.URL.Query().Get("offset")))
		items, e := h.Admissions.List(r.Context(), pid, r.URL.Query().Get("status"), q.Limit, q.Offset)
		if e != nil {
			writeErr(w, e)
			return
		}
		write(w, 200, items)
		return
	}
	write(w, 405, nil)
}
func atoi(v string) int { n, _ := strconv.Atoi(v); return n }
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
func writeErr(w http.ResponseWriter, e error) {
	code := 500
	if errors.Is(e, domain.ErrForbidden) {
		code = 403
	} else if errors.Is(e, domain.ErrNotFound) {
		code = 404
	} else if errors.Is(e, domain.ErrConflict) || errors.Is(e, domain.ErrCapacity) {
		code = 409
	} else if errors.Is(e, domain.ErrInvalidState) {
		code = 422
	}
	write(w, code, map[string]string{"error": e.Error()})
}
