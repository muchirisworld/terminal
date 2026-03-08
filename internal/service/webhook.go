package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/muchirisworld/terminal/internal/logger"
	"github.com/muchirisworld/terminal/internal/store"
)

// Event envelope parsing
type Event struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// WebhookService processes incoming webhooks
type WebhookService struct {
	store  *store.Store
	logger *slog.Logger
}

// NewWebhookService creates a new WebhookService
func NewWebhookService(s *store.Store, log *slog.Logger) *WebhookService {
	return &WebhookService{
		store:  s,
		logger: log,
	}
}

// Process handles a single webhook event
func (s *WebhookService) Process(ctx context.Context, eventID string, rawBody []byte) error {
	var event Event
	if err := json.Unmarshal(rawBody, &event); err != nil {
		return fmt.Errorf("failed to parse event envelope: %w", err)
	}

	logger.Add(ctx, "webhook_type", event.Type)

	// Try to insert the event
	alreadyProcessed, err := s.store.InsertWebhookEvent(ctx, eventID, "clerk", event.Type, rawBody)
	if err != nil {
		return fmt.Errorf("failed to insert webhook event: %w", err)
	}
	if alreadyProcessed {
		logger.Add(ctx, "webhook_status", "already_processed")
		return nil
	}

	handler, exists := GetClerkHandler(event.Type)
	if !exists {
		// Mark as processed if unsupported
		logger.Add(ctx, "webhook_status", "unsupported_type")
		if err := s.store.MarkWebhookEventProcessed(ctx, eventID); err != nil {
			return fmt.Errorf("failed to update unsupported event: %w", err)
		}
		return nil
	}

	// Run the specific projection handler in a transaction
	err = s.store.ExecTx(ctx, func(txStore *store.Store) error {
		return handler(ctx, txStore, event)
	})

	if err != nil {
		logger.Add(ctx, "webhook_status", "handler_failed")
		logger.Add(ctx, "error", err.Error())

		// Create a separate connection/transaction to record the error
		if updateErr := s.store.UpdateWebhookEventError(context.Background(), eventID, err.Error()); updateErr != nil {
			// This is background context, so we still log it using the injected logger to be safe
			s.logger.Error("failed to update webhook event error state", "id", eventID, "err", updateErr)
		}

		// Return the error so the caller can return a 500 status and trigger retry
		return err
	}

	// Handler succeeded
	if err := s.store.MarkWebhookEventProcessed(ctx, eventID); err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}

	logger.Add(ctx, "webhook_status", "processed")
	return nil
}
