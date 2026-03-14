package app

import (
	nethttp "net/http"

	"github.com/chgc/golf_team_manager/backend/internal/config"
	apihttp "github.com/chgc/golf_team_manager/backend/internal/http"
)

type App struct {
	server *nethttp.Server
}

func New(cfg config.Config) *App {
	return &App{
		server: &nethttp.Server{
			Addr:              cfg.HTTP.Address(),
			Handler:           apihttp.NewRouter(),
			ReadHeaderTimeout: cfg.HTTP.ReadTimeout,
		},
	}
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}
