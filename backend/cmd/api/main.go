package main

import (
	"errors"
	"log"
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/app"
	"github.com/chgc/golf_team_manager/backend/internal/config"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	if err := application.Run(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
		log.Fatalf("run app: %v", err)
	}
}
