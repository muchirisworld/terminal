-- +goose Up
-- +goose StatementBegin
CREATE TYPE org_member_role AS ENUM (
    'admin',
    'member'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE memberships (
    org_id varchar NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id varchar NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    role org_member_role NOT NULL DEFAULT 'member',

    is_active boolean NOT NULL DEFAULT true,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    -- Composite Primary Key
    PRIMARY KEY (user_id, org_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_memberships_org_id ON memberships(org_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_memberships_updated_at
BEFORE UPDATE ON "memberships"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS memberships;
DROP TYPE IF EXISTS org_member_role;
-- +goose StatementEnd
