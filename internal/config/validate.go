package config

import (
	"fmt"
	"net"
)

func (c Config) Validate() error {
	if c.HTTPAddr == "" || c.DBPath == "" {
		return fmt.Errorf("address and database are required")
	}
	if _, e := net.ResolveTCPAddr("tcp", c.HTTPAddr); e != nil {
		return fmt.Errorf("invalid http address: %w", e)
	}
	if c.SessionTTLSeconds < 60 {
		return fmt.Errorf("session ttl too short")
	}
	return nil
}
