-- +goose Up
-- +goose StatementBegin
CREATE TABLE organizations (
    id varchar PRIMARY KEY NOT NULL,

    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    image_url text,
    metadata jsonb DEFAULT '{}'::jsonb,
    max_allowed_memberships INTEGER NOT NULL DEFAULT 5,

    created_by varchar REFERENCES users(id) ON DELETE SET NULL,

    is_active boolean NOT NULL DEFAULT true,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_organizations_updated_at
BEFORE UPDATE ON "organizations"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_organizations_created_by ON organizations(created_by);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd
