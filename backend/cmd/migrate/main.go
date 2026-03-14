package main

import (
	"context"
	"log"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := appdb.Open(cfg.DB)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer database.Close()

	if err := appdb.RunMigrations(context.Background(), database); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
}
