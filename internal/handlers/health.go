package handlers

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
)

// HealthHandler is the handler for the healthz endpoint.
type HealthHandler struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *sqlx.DB, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		DB:     db,
		Logger: logger,
	}
}

// Healthz is the handler for the healthz endpoint.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Readyz is the handler for the readyz endpoint.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.PingContext(r.Context()); err != nil {
		h.Logger.Error("database not ready", "err", err)
		http.Error(w, "database not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
