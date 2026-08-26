package worker

import (
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	processed atomic.Int64
	failed    atomic.Int64
	mu        sync.Mutex
	last      time.Time
}

func (m *Metrics) Processed() { m.processed.Add(1); m.mu.Lock(); m.last = time.Now(); m.mu.Unlock() }
func (m *Metrics) Failed()    { m.failed.Add(1); m.mu.Lock(); m.last = time.Now(); m.mu.Unlock() }
func (m *Metrics) Snapshot() (int64, int64, time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.processed.Load(), m.failed.Load(), m.last
}
