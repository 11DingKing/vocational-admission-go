package config

import "testing"

func TestDefaultsAndValidation(t *testing.T) {
	c := Config{HTTPAddr: ":8080", DBPath: "admission.db", SessionTTLSeconds: 3600, WorkerIntervalSeconds: 5}
	if c.Validate() != nil {
		t.Fatal("valid config rejected")
	}
	if c.WorkerInterval() != 5 || c.SessionTTL() != 3600 {
		t.Fatal("defaults")
	}
	if len(DefaultHeaders()) != 2 {
		t.Fatal("headers")
	}
	if len(SupportedTimezones()) != 3 {
		t.Fatal("zones")
	}
	if !IsProduction(c) {
		t.Fatal("production")
	}
	for _, ttl := range []int{60, 120, 300, 600, 900, 1800, 3600, 7200, 86400} {
		c.SessionTTLSeconds = ttl
		if c.Validate() != nil {
			t.Errorf("ttl %d", ttl)
		}
	}
}
