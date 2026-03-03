package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/muchirisworld/terminal/internal/handlers"
	imiddleware "github.com/muchirisworld/terminal/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(logger *slog.Logger, healthRouter, userRouter, webhookRouter http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(imiddleware.Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/health", healthRouter)
	r.Mount("/users", userRouter)
	r.Mount("/webhooks", webhookRouter)

	return r
}

func RegisterWebhookRoutes(h *handlers.WebhookHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/clerk", h.HandleClerk)
	return r
}

func RegisterHealthRoutes(h *handlers.HealthHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", h.Healthz)
	r.Get("/readyz", h.Readyz)
	return r
}

func RegisterUserRoutes(h *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	return r
}
