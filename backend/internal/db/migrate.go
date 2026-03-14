package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"

	"github.com/chgc/golf_team_manager/backend/migrations"
)

const createSchemaMigrationsTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func RunMigrations(ctx context.Context, database *sql.DB) error {
	migrationFiles, err := fs.Glob(migrations.Files, "*.sql")
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}

	sort.Strings(migrationFiles)

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, createSchemaMigrationsTableSQL); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	appliedVersions, err := loadAppliedVersions(ctx, tx)
	if err != nil {
		return err
	}

	for _, migrationFile := range migrationFiles {
		if _, ok := appliedVersions[migrationFile]; ok {
			continue
		}

		migrationSQL, err := migrations.Files.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", migrationFile, err)
		}

		if _, err := tx.ExecContext(ctx, string(migrationSQL)); err != nil {
			return fmt.Errorf("apply migration %s: %w", migrationFile, err)
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO schema_migrations (version) VALUES (?)`,
			migrationFile,
		); err != nil {
			return fmt.Errorf("record migration %s: %w", migrationFile, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migrations: %w", err)
	}

	return nil
}

func loadAppliedVersions(ctx context.Context, tx *sql.Tx) (map[string]struct{}, error) {
	rows, err := tx.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("query applied migrations: %w", err)
	}
	defer rows.Close()

	appliedVersions := make(map[string]struct{})
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}

		appliedVersions[version] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return appliedVersions, nil
}
