
-- +goose Up
ALTER TABLE auth.users
ADD COLUMN provider TEXT DEFAULT 'email';

-- +goose Down
ALTER TABLE auth.users
DROP COLUMN provider;

