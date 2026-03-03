package store

import (
	"context"

	"github.com/lib/pq"
)

// InsertWebhookEvent tries to insert a new webhook event.
// Returns a boolean indicating if it was already processed (true) and should be skipped.
func (s *Store) InsertWebhookEvent(ctx context.Context, id, provider, eventType string, payload []byte) (bool, error) {

	_, err := s.dbtx.ExecContext(ctx, `
		INSERT INTO webhook_events (id, provider, type, payload, attempts)
		VALUES ($1, $2, $3, $4, 1)
	`, id, provider, eventType, payload)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // unique_violation
			// Try to acquire lock on the row for processing
			var processed bool
			errCheck := s.dbtx.QueryRowContext(ctx, "SELECT processed_at IS NOT NULL FROM webhook_events WHERE id = $1 FOR UPDATE SKIP LOCKED", id).Scan(&processed)
			if errCheck == sql.ErrNoRows {
				// Another process has the lock, skip this
				return true, nil
			}
			if errCheck != nil {
				return false, errCheck
			}
			if processed {
				return true, nil // Already processed, skip
			}
			return false, nil // Not processed, allow retry
		}
		return false, err
	}

	return false, nil
}


// MarkWebhookEventProcessed marks an event as processed.
func (s *Store) MarkWebhookEventProcessed(ctx context.Context, id string) error {
	_, err := s.dbtx.ExecContext(ctx, `UPDATE webhook_events SET processed_at=now() WHERE id=$1`, id)
	return err
}

// UpdateWebhookEventError updates the error state and increments attempts.
func (s *Store) UpdateWebhookEventError(ctx context.Context, id, errStr string) error {
	// Notice we might need a separate connection/transaction for this if the current tx is aborted.
	// Typically, we use the bare s.db for this to bypass the broken transaction.
	_, err := s.db.ExecContext(ctx, `
		UPDATE webhook_events
		SET last_error=$1, attempts=attempts+1
		WHERE id=$2
	`, errStr, id)
	return err
}
