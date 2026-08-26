package worker

import (
	"context"
	"errors"
	"github.com/11DingKing/vocational-admission-go/internal/repository"
	"log/slog"
	"time"
)

type Worker struct {
	Jobs     repository.JobRepo
	Interval time.Duration
	Stop     chan struct{}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.Stop:
			return
		case <-ticker.C:
			w.once(ctx)
		}
	}
}
func (w *Worker) once(ctx context.Context) {
	j, e := w.Jobs.Claim(ctx)
	if errors.Is(e, context.Canceled) {
		return
	}
	if e != nil {
		return
	}
	if e = w.handle(ctx, j); e != nil {
		_ = w.Jobs.Fail(ctx, j.ID, e.Error())
		slog.Error("job failed", "id", j.ID, "error", e)
		return
	}
	_ = w.Jobs.Complete(ctx, j.ID)
}
func (w *Worker) handle(ctx context.Context, j interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
