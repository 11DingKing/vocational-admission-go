package worker

import (
	"context"
	"time"
)

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 8 {
		attempt = 8
	}
	return p.BaseDelay * time.Duration(1<<(attempt-1))
}
func Retry(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	var err error
	for n := 1; n <= p.MaxAttempts; n++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		d := p.Delay(n)
		t := time.NewTimer(d)
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
	return err
}
