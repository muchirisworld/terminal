-- +goose Up
-- +goose StatementBegin
CREATE TABLE webhook_events (
    instance_id text PRIMARY KEY,
    provider text NOT NULL,
    type text NOT NULL,
    payload jsonb NOT NULL,
    received_at timestamptz NOT NULL DEFAULT now(),
    processed_at timestamptz,
    attempts integer NOT NULL DEFAULT 0,
    last_error text
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_webhook_events_unprocessed ON webhook_events(processed_at) WHERE processed_at IS NULL;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_webhook_events_provider ON webhook_events(provider);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_webhook_events_provider;
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_webhook_events_unprocessed;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS webhook_events;
-- +goose StatementEnd
