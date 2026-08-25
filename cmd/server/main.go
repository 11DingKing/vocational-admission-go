package main

import (
	"context"
	"errors"
	"github.com/11DingKing/vocational-admission-go/internal/audit"
	"github.com/11DingKing/vocational-admission-go/internal/auth"
	"github.com/11DingKing/vocational-admission-go/internal/config"
	"github.com/11DingKing/vocational-admission-go/internal/httpapi"
	"github.com/11DingKing/vocational-admission-go/internal/middleware"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"github.com/11DingKing/vocational-admission-go/internal/service"
	"github.com/11DingKing/vocational-admission-go/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	cfg := config.Load()
	ctx := context.Background()
	db, e := storage.Open(ctx, cfg.DBPath)
	if e != nil {
		return e
	}
	defer db.Close()
	if e = storage.Migrate(ctx, db.SQL); e != nil {
		return e
	}
	users := repository.UserRepo{DB: db.SQL}
	sessions := repository.SessionRepo{DB: db.SQL}
	a := auth.Service{Users: users, Sessions: sessions, TTL: time.Duration(cfg.SessionTTLSeconds) * time.Second}
	plans := repository.PlanRepo{DB: db.SQL}
	apps := repository.ApplicationRepo{DB: db.SQL}
	svc := service.AdmissionService{DB: db, Plans: plans, Apps: apps, Audits: audit.Logger{Repo: repository.AuditRepo{DB: db.SQL}}, Jobs: repository.JobRepo{DB: db.SQL}}
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: middleware.Recover(middleware.RequestID(httpapi.New(a, svc).Mux)), ReadHeaderTimeout: 5 * time.Second}
	stop, stopCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stopCancel()
	go func() { <-stop.Done(); _ = srv.Shutdown(context.Background()) }()
	slog.Info("server listening", "addr", cfg.HTTPAddr)
	e = srv.ListenAndServe()
	if errors.Is(e, http.ErrServerClosed) {
		return nil
	}
	return e
}
