package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/logger"
	"github.com/muchirisworld/terminal/internal/service"
	svix "github.com/svix/svix-webhooks/go"
)

type WebhookHandler struct {
	svc    *service.WebhookService
	secret string
	logger *slog.Logger
}

func NewWebhookHandler(svc *service.WebhookService, cfg *config.Config, log *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		svc:    svc,
		secret: cfg.ClerkWebhookSecret,
		logger: log,
	}
}

func (h *WebhookHandler) HandleClerk(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Add(r.Context(), "error", err.Error())
		logger.Add(r.Context(), "webhook_error", "failed to read webhook body")
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// Verify signature if secret is provided
	if h.secret != "" {
		wh, err := svix.NewWebhook(h.secret)
		if err != nil {
			logger.Add(r.Context(), "error", err.Error())
			logger.Add(r.Context(), "webhook_error", "failed to create webhook verifier")
			http.Error(w, "internal configuration error", http.StatusInternalServerError)
			return
		}

		err = wh.Verify(payload, r.Header)
		if err != nil {
			logger.Add(r.Context(), "error", err.Error())
			logger.Add(r.Context(), "webhook_error", "invalid webhook signature")
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	} else {
		logger.Add(r.Context(), "webhook_warning", "webhook signature verification skipped because secret is not configured")
	}

	eventID := r.Header.Get("svix-id")
	if eventID == "" {
		hash := sha256.Sum256(payload)
		eventID = hex.EncodeToString(hash[:])
	}
	logger.Add(r.Context(), "webhook_id", eventID)

	err = h.svc.Process(r.Context(), eventID, payload)
	if err != nil {
		logger.Add(r.Context(), "error", err.Error())
		http.Error(w, "failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
