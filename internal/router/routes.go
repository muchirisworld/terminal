package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/handlers"
	imiddleware "github.com/muchirisworld/terminal/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(cfg *config.Config, logger *slog.Logger, healthRouter, userRouter, webhookRouter, catalogRouter, inventoryRouter http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(imiddleware.Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/health", healthRouter)
	r.Mount("/users", userRouter)
	r.Mount("/webhooks", webhookRouter)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(imiddleware.AuthMiddleware(cfg, logger))
		r.Mount("/catalog", catalogRouter)
		r.Mount("/inventory", inventoryRouter)
	})

	return r
}

func RegisterCatalogRoutes(h *handlers.CatalogHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/products", h.CreateProduct)
	r.Get("/products", h.ListProducts)
	r.Get("/products/{productID}", h.GetProduct)
	r.Patch("/products/{productID}", h.UpdateProduct)
	r.Post("/products/{productID}/archive", h.ArchiveProduct)

	r.Post("/products/{productID}/variants", h.CreateVariant)
	r.Get("/products/{productID}/variants", h.ListVariantsByProduct)
	r.Get("/variants/{variantID}", h.GetVariant)
	r.Patch("/variants/{variantID}", h.UpdateVariant)

	return r
}

func RegisterInventoryRoutes(h *handlers.InventoryHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/products/{productID}/conversions", h.UpsertConversion)
	r.Get("/products/{productID}/conversions", h.ListConversionsByProduct)

	r.Post("/variants/{variantID}/receipt", h.CreateReceipt)
	r.Post("/variants/{variantID}/adjustment", h.CreateAdjustment)
	r.Post("/variants/{variantID}/reserve", h.ReserveInventory)
	r.Get("/variants/{variantID}/stock", h.GetVariantStock)

	r.Post("/reservations/{reservationID}/release", h.ReleaseReservation)

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
