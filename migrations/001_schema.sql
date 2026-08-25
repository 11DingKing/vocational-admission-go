-- Canonical schema is embedded in internal/storage/migrate.go for single-binary deployments.
-- This file documents migration 001 and is kept in sync with the executable migration.
CREATE TABLE schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);
