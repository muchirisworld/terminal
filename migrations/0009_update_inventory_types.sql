-- +goose Up
-- +goose StatementBegin
ALTER TABLE inventory_events 
    ALTER COLUMN quantity_change TYPE numeric(18,4);

ALTER TABLE inventory_reservations 
    ALTER COLUMN quantity TYPE numeric(18,4);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE inventory_events 
    ALTER COLUMN quantity_change TYPE bigint;

ALTER TABLE inventory_reservations 
    ALTER COLUMN quantity TYPE bigint;
-- +goose StatementEnd
