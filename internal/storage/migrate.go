package storage

import (
	"context"
	"database/sql"
	"fmt"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);`,
	`CREATE TABLE IF NOT EXISTS users(id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL, role TEXT NOT NULL, active INTEGER NOT NULL DEFAULT 1, created_at TEXT NOT NULL); CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY, user_id INTEGER NOT NULL REFERENCES users(id), expires_at TEXT NOT NULL, revoked_at TEXT, created_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_id);`,
	`CREATE TABLE IF NOT EXISTS province_rules(id INTEGER PRIMARY KEY AUTOINCREMENT, province TEXT NOT NULL, batch TEXT NOT NULL, min_score INTEGER NOT NULL, rank_limit INTEGER NOT NULL, allow_transfer INTEGER NOT NULL, version INTEGER NOT NULL DEFAULT 1, updated_at TEXT NOT NULL, UNIQUE(province,batch)); CREATE TABLE IF NOT EXISTS plans(id INTEGER PRIMARY KEY AUTOINCREMENT, year INTEGER NOT NULL, province TEXT NOT NULL, name TEXT NOT NULL, status TEXT NOT NULL, total_capacity INTEGER NOT NULL, used_capacity INTEGER NOT NULL DEFAULT 0, version INTEGER NOT NULL DEFAULT 1, locked_at TEXT, created_at TEXT NOT NULL, UNIQUE(year,province)); CREATE INDEX IF NOT EXISTS idx_plans_status ON plans(status);`,
	`CREATE TABLE IF NOT EXISTS major_groups(id INTEGER PRIMARY KEY AUTOINCREMENT, plan_id INTEGER NOT NULL REFERENCES plans(id), code TEXT NOT NULL, name TEXT NOT NULL, capacity INTEGER NOT NULL, used_capacity INTEGER NOT NULL DEFAULT 0, version INTEGER NOT NULL DEFAULT 1, UNIQUE(plan_id,code)); CREATE TABLE IF NOT EXISTS applications(id INTEGER PRIMARY KEY AUTOINCREMENT, plan_id INTEGER NOT NULL REFERENCES plans(id), major_group_id INTEGER NOT NULL REFERENCES major_groups(id), student_no TEXT NOT NULL, score INTEGER NOT NULL, rank INTEGER NOT NULL, preferences TEXT NOT NULL, status TEXT NOT NULL, transfer_accepted INTEGER NOT NULL DEFAULT 0, idempotency_key TEXT NOT NULL UNIQUE, version INTEGER NOT NULL DEFAULT 1, submitted_at TEXT NOT NULL, updated_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_apps_plan_status ON applications(plan_id,status);`,
	`CREATE TABLE IF NOT EXISTS decisions(id INTEGER PRIMARY KEY AUTOINCREMENT, application_id INTEGER NOT NULL REFERENCES applications(id), actor_id INTEGER NOT NULL REFERENCES users(id), from_status TEXT NOT NULL, to_status TEXT NOT NULL, reason TEXT NOT NULL, request_id TEXT NOT NULL, created_at TEXT NOT NULL); CREATE TABLE IF NOT EXISTS audit_events(id INTEGER PRIMARY KEY AUTOINCREMENT, actor_id INTEGER REFERENCES users(id), entity TEXT NOT NULL, entity_id INTEGER NOT NULL, action TEXT NOT NULL, outcome TEXT NOT NULL, request_id TEXT NOT NULL, details TEXT NOT NULL, created_at TEXT NOT NULL); CREATE TABLE IF NOT EXISTS jobs(id INTEGER PRIMARY KEY AUTOINCREMENT, kind TEXT NOT NULL, entity_id INTEGER NOT NULL, attempts INTEGER NOT NULL DEFAULT 0, available_at TEXT NOT NULL, completed_at TEXT, last_error TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_jobs_ready ON jobs(completed_at,available_at);`,
}

func Migrate(ctx context.Context, db *sql.DB) error {
	for i, m := range migrations {
		if err := apply(ctx, db, i+1, m); err != nil {
			return fmt.Errorf("migration %d: %w", i+1, err)
		}
	}
	return nil
}
func apply(ctx context.Context, db *sql.DB, v int, sqlText string) error {
	if v == 1 {
		if _, err := db.ExecContext(ctx, sqlText); err != nil {
			return err
		}
		_, err := db.ExecContext(ctx, "INSERT OR IGNORE INTO schema_migrations(version,applied_at) VALUES(?,datetime('now'))", v)
		return err
	}
	var n int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(1) FROM schema_migrations WHERE version=?", v).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, sqlText); err == nil {
		_, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations(version,applied_at) VALUES(?,datetime('now'))", v)
	}
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
