package config

import "os"

type Config struct {
	HTTPAddr              string
	DBPath                string
	SessionTTLSeconds     int
	WorkerIntervalSeconds int
}

func Load() Config {
	c := Config{HTTPAddr: ":8080", DBPath: "admission.db", SessionTTLSeconds: 3600, WorkerIntervalSeconds: 5}
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		c.DBPath = v
	}
	return c
}
