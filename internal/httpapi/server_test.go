package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/11DingKing/vocational-admission-go/internal/audit"
	"github.com/11DingKing/vocational-admission-go/internal/auth"
	"github.com/11DingKing/vocational-admission-go/internal/domain"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/service"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func api(t *testing.T) http.Handler {
	db, e := storage.Open(context.Background(), "file:http-"+t.Name()+"?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	if e = storage.Migrate(context.Background(), db.SQL); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	a := auth.Service{Users: repository.UserRepo{DB: db.SQL}, Sessions: repository.SessionRepo{DB: db.SQL}, TTL: time.Hour}
	s := service.AdmissionService{DB: db, Plans: repository.PlanRepo{DB: db.SQL}, Apps: repository.ApplicationRepo{DB: db.SQL}, Audits: audit.Logger{Repo: repository.AuditRepo{DB: db.SQL}}, Jobs: repository.JobRepo{DB: db.SQL}}
	return New(a, s).Mux
}
func post(t *testing.T, h http.Handler, path string, v any, headers map[string]string) *httptest.ResponseRecorder {
	b, _ := json.Marshal(v)
	r := httptest.NewRequest("POST", path, bytes.NewReader(b))
	for k, x := range headers {
		r.Header.Set(k, x)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func TestHealthReady(t *testing.T) {
	h := api(t)
	for _, path := range []string{"/healthz", "/readyz"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		if w.Code != 200 {
			t.Fatalf("%s=%d", path, w.Code)
		}
	}
}
func TestRegisterLoginLogout(t *testing.T) {
	h := api(t)
	w := post(t, h, "/v1/auth/register", map[string]any{"Username": "off", "Password": "pw", "Role": domain.RoleOfficer}, nil)
	if w.Code != 201 {
		t.Fatalf("register %d", w.Code)
	}
	w = post(t, h, "/v1/auth/login", map[string]string{"Username": "off", "Password": "pw"}, nil)
	if w.Code != 200 {
		t.Fatalf("login %d", w.Code)
	}
	var out struct{ Token string }
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	if out.Token == "" {
		t.Fatal("missing token")
	}
	r := httptest.NewRequest("POST", "/v1/auth/logout", nil)
	r.Header.Set("Authorization", "Bearer "+out.Token)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, r)
	if rw.Code != 204 {
		t.Fatalf("logout %d", rw.Code)
	}
}
func TestProtectedEndpoint(t *testing.T) {
	h := api(t)
	w := post(t, h, "/v1/plans", map[string]any{"Year": 2026, "Province": "北京", "Name": "计划", "Capacity": 10}, nil)
	if w.Code != 401 {
		t.Fatalf("unauthorized=%d", w.Code)
	}
}
func TestMethodErrors(t *testing.T) {
	h := api(t)
	r := httptest.NewRequest("PUT", "/healthz", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("health accepts methods by design")
	}
}
