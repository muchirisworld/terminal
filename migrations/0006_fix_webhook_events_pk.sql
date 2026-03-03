-- +goose Up
-- +goose StatementBegin
ALTER TABLE webhook_events RENAME COLUMN instance_id TO id;
DELETE FROM webhook_events; -- Clear duplicates that were stuck
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE webhook_events RENAME COLUMN id TO instance_id;
-- +goose StatementEnd
