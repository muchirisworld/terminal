package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/service"
	svix "github.com/svix/svix-webhooks/go"
)

type WebhookHandler struct {
	svc    *service.WebhookService
	secret string
	logger *slog.Logger
}

func NewWebhookHandler(svc *service.WebhookService, cfg *config.Config, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		svc:    svc,
		secret: cfg.ClerkWebhookSecret,
		logger: logger,
	}
}

func (h *WebhookHandler) HandleClerk(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read webhook body", "error", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	// Verify signature if secret is provided
	if h.secret != "" {
		wh, err := svix.NewWebhook(h.secret)
		if err != nil {
			h.logger.Error("failed to create webhook verifier", "error", err)
			http.Error(w, "internal configuration error", http.StatusInternalServerError)
			return
		}

		err = wh.Verify(payload, r.Header)
		if err != nil {
			h.logger.Warn("invalid webhook signature", "error", err)
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	} else {
		h.logger.Warn("webhook signature verification skipped because secret is not configured")
	}

	eventID := r.Header.Get("svix-id")
	if eventID == "" {
		hash := sha256.Sum256(payload)
		eventID = hex.EncodeToString(hash[:])
	}

	err = h.svc.Process(r.Context(), eventID, payload)
	if err != nil {
		// The service already logged the error, so we just return 500
		http.Error(w, "failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
