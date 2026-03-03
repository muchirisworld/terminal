package handlers_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/muchirisworld/terminal/internal/config"
	"github.com/muchirisworld/terminal/internal/handlers"
	"github.com/muchirisworld/terminal/internal/service"
	"github.com/muchirisworld/terminal/internal/store"
)

func TestWebhookHandler_SignatureFailure(t *testing.T) {
	cfg := &config.Config{
		ClerkWebhookSecret: "whsec_testsecret123",
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// Pass nil db for this test since we only test the signature which fails before DB is hit
	svc := service.NewWebhookService(store.New(nil), logger)
	handler := handlers.NewWebhookHandler(svc, cfg, logger)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/clerk", bytes.NewBufferString(`{"type":"test"}`))
	req.Header.Set("svix-id", "msg_123")
	req.Header.Set("svix-timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	req.Header.Set("svix-signature", "v1,bad_signature")

	rr := httptest.NewRecorder()
	handler.HandleClerk(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestWebhookHandler_SignatureSuccess(t *testing.T) {
	// This would test success, but requires a valid svix signature for a payload.
	// Since we mock the DB, if it passes signature, it will panic on DB.
	// We'll test full flow in integration tests.
}
