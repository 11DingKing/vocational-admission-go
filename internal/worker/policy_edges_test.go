package worker

import "testing"

func TestRetryPolicyEdges(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 4, BaseDelay: 1}
	if p.Delay(0) != 1 {
		t.Fatal("zero attempt")
	}
	if p.Delay(20) != 128 {
		t.Fatal("capped delay")
	}
	m := Metrics{}
	for i := 0; i < 10; i++ {
		m.Processed()
	}
	for i := 0; i < 3; i++ {
		m.Failed()
	}
	a, b, _ := m.Snapshot()
	if a != 10 || b != 3 {
		t.Fatalf("%d %d", a, b)
	}
}
