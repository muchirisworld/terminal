-- +goose Up
-- +goose StatementBegin
CREATE TYPE invitation_status AS ENUM (
    'pending',
    'accepted',
    'revoked',
    'expired'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE invitations (
    id varchar PRIMARY KEY NOT NULL,

    org_id varchar NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    email email_citext NOT NULL,

    invited_by_user_id varchar REFERENCES users(id) ON DELETE SET NULL,
    accepted_by_user_id varchar REFERENCES users(id) ON DELETE SET NULL,

    role org_member_role NOT NULL DEFAULT 'member',

    status invitation_status NOT NULL DEFAULT 'pending',

    expires_at timestamptz,
    accepted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_pending_invite
ON invitations(org_id, email)
WHERE status = 'pending';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_invitations_updated_at
BEFORE UPDATE ON "invitations"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_invitations_org_id ON invitations(org_id);
CREATE INDEX idx_invitations_email ON invitations(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invitations;
DROP TYPE IF EXISTS invitation_status;
-- +goose StatementEnd
