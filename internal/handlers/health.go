package handlers

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/muchirisworld/terminal/internal/logger"
)

// HealthHandler is the handler for the healthz endpoint.
type HealthHandler struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *sqlx.DB, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		DB:     db,
		Logger: log,
	}
}

// Healthz is the handler for the healthz endpoint.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz is the handler for the readyz endpoint.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.PingContext(r.Context()); err != nil {
		logger.Add(r.Context(), "error", err.Error())
		logger.Add(r.Context(), "health_error", "database not ready")
		http.Error(w, "database not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
