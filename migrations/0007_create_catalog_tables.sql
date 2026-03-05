-- +goose Up
-- +goose StatementBegin
CREATE TYPE product_status AS ENUM (
    'active',
    'archived'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE products (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    name text NOT NULL,
    description text,
    base_unit text NOT NULL,
    status product_status NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_products_organization_id ON products(organization_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_products_updated_at
BEFORE UPDATE ON "products"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE product_variants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku text NOT NULL,
    barcode text,
    price numeric(12,2) NOT NULL,
    cost numeric(12,2),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE(organization_id, sku)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_variants_org_product ON product_variants(organization_id, product_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_product_variants_updated_at
BEFORE UPDATE ON "product_variants"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE product_images (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id varchar NOT NULL,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_key text NOT NULL,
    position int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_product_images_org_product ON product_images(organization_id, product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS product_status;
-- +goose StatementEnd
