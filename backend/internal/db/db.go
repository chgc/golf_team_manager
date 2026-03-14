package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	_ "modernc.org/sqlite"
)

func Open(cfg config.DBConfig) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	database, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	database.SetMaxOpenConns(1)

	if _, err := database.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		database.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return database, nil
}
