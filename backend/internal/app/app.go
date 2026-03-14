package app

import (
	"context"
	"database/sql"
	"fmt"
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	appdb "github.com/chgc/golf_team_manager/backend/internal/db"
	apihttp "github.com/chgc/golf_team_manager/backend/internal/http"
)

type App struct {
	database *sql.DB
	server   *nethttp.Server
}

func New(cfg config.Config) (*App, error) {
	database, err := appdb.Open(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if cfg.DB.AutoMigrate {
		if err := appdb.RunMigrations(context.Background(), database); err != nil {
			database.Close()
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	return &App{
		database: database,
		server: &nethttp.Server{
			Addr:              cfg.HTTP.Address(),
			Handler:           apihttp.NewRouter(database),
			ReadHeaderTimeout: cfg.HTTP.ReadTimeout,
		},
	}, nil
}

func (a *App) Run() error {
	defer a.database.Close()

	return a.server.ListenAndServe()
}
