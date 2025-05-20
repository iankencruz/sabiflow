
-- +goose Up
ALTER TABLE auth.sessions
DROP COLUMN IF EXISTS data,
DROP COLUMN IF EXISTS expiry;

-- +goose Down
ALTER TABLE auth.sessions
ADD COLUMN data BYTEA NOT NULL DEFAULT ''::bytea,
ADD COLUMN expiry TIMESTAMPTZ NOT NULL DEFAULT now() + interval '7 days';

