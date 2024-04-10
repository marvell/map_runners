package main

import (
	"context"
	"database/sql"
	"fmt"
)

type Application struct {
	cfg    *Config
	db     *sql.DB
	store  *Store
	strava *Strava
}

func NewApplication(cfg *Config) *Application {
	app := &Application{
		cfg: cfg,
	}

	app.InitStrava()

	return app
}

func (app *Application) Run(ctx context.Context) error {
	if err := app.InitPostgres(ctx); err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}
	app.store = NewStore(app.db)

	if err := app.RunHTTPServer(ctx); err != nil {
		return fmt.Errorf("run http-server: %w", err)
	}

	return nil
}

func (app *Application) Close(ctx context.Context) error {
	app.ClosePostgres()
	return nil
}
