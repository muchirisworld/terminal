-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'email_citext') THEN
        CREATE DOMAIN email_citext AS citext
            CHECK (VALUE ~* '^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,63}$');
    END IF;
END$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trg_set_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END$$;
-- +goose StatementEnd

CREATE TABLE users (
    "id" varchar PRIMARY KEY NOT NULL,
    "name" text NOT NULL,
    "email" email_citext NOT NULL UNIQUE,
    "email_verified" boolean DEFAULT false NOT NULL,
    "image" text,
    "created_at" timestamptz DEFAULT now() NOT NULL,
    "updated_at" timestamptz DEFAULT now() NOT NULL
);

-- +goose StatementBegin
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON "users"
FOR EACH ROW EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS "users";
DROP DOMAIN IF EXISTS email_citext;
DROP FUNCTION IF EXISTS trg_set_updated_at;