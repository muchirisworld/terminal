package db

import (
	"context"
	"time"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver
	"github.com/muchirisworld/terminal/internal/config"
)

// New creates a new database connection.
func New(cfg *config.Config, ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
