package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (app *Application) InitPostgres(ctx context.Context) error {
	db, err := sql.Open("postgres", app.cfg.PostgresDSN)
	if err != nil {
		return fmt.Errorf("connection to postgres: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}
	app.db = db

	return nil
}

func (app *Application) ClosePostgres() {
	app.db.Close()
}
