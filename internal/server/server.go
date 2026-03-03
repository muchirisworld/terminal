package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/router"
)

// New creates a new HTTP server.
func New(cfg *config.Config, logger *slog.Logger, healthHandler, userHandler, webhookHandler http.Handler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: router.RegisterRoutes(logger, healthHandler, userHandler, webhookHandler),
	}
}
