-- +goose Up
-- +goose StatementBegin
CREATE TYPE inventory_event_type AS ENUM (
    'purchase_receipt',
    'order_fulfillment',
    'adjustment',
    'return',
    'case_break'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE inventory_source_type AS ENUM (
    'receipt',
    'order',
    'manual',
    'system'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE inventory_reservation_status AS ENUM (
    'active',
    'released',
    'expired'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE unit_conversions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    unit_from text NOT NULL,
    unit_to text NOT NULL,
    factor numeric(18,6) NOT NULL,
    precision int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE(organization_id, product_id, unit_from, unit_to)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_unit_conversions_updated_at
BEFORE UPDATE ON "unit_conversions"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE inventory_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    product_variant_id uuid NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    event_type inventory_event_type NOT NULL,
    quantity_change bigint NOT NULL,
    source_type inventory_source_type,
    source_id uuid,
    note text,
    created_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_inventory_events_lookup ON inventory_events(organization_id, product_variant_id, created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE inventory_reservations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    product_variant_id uuid NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    order_id uuid,
    quantity bigint NOT NULL,
    status inventory_reservation_status NOT NULL DEFAULT 'active',
    expires_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    released_at timestamptz
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_inventory_reservations_lookup ON inventory_reservations(organization_id, product_variant_id, status, expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventory_reservations;
DROP TABLE IF EXISTS inventory_events;
DROP TABLE IF EXISTS unit_conversions;
DROP TYPE IF EXISTS inventory_reservation_status;
DROP TYPE IF EXISTS inventory_source_type;
DROP TYPE IF EXISTS inventory_event_type;
-- +goose StatementEnd
