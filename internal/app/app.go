package app

import "context"

interface {
	Printf(format string, v ...any)
}

type Config struct {
}

type App struct {
	cfg Config
}

func New(cfg Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) error {

}
