package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/db"
	"github.com/muchirisworld/terminal/internal/handlers"
	"github.com/muchirisworld/terminal/internal/logger"
	"github.com/muchirisworld/terminal/internal/router"
	"github.com/muchirisworld/terminal/internal/server"
	"github.com/muchirisworld/terminal/internal/service"
	"github.com/muchirisworld/terminal/internal/store"

	"github.com/jmoiron/sqlx"
)

// App is the application container.
type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *sqlx.DB
	Server *http.Server
}

// New creates a new App.
func New(ctx context.Context) (*App, error) {
	cfg := config.New()

	log := logger.New(cfg)

	database, err := db.New(cfg, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	healthHandler := handlers.NewHealthHandler(database, log)
	healthRouter := router.RegisterHealthRoutes(healthHandler)

	userStore := store.NewUserStore(database)
	userService := service.NewUserService(userStore)
	userHandler := handlers.NewUserHandler(userService, log)
	userRouter := router.RegisterUserRoutes(userHandler)

	return &App{
		Config: cfg,
		Logger: log,
		DB:     database,
		Server: server.New(cfg, log, healthRouter, userRouter),
	}, nil
}

// Run starts the application.
func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("starting server", "port", a.Config.HTTPPort)
	return a.Server.ListenAndServe()
}

// Close gracefully shuts down the application.
func (a *App) Close(ctx context.Context) error {
	a.Logger.Info("shutting down server")

	if err := a.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	if err := a.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
