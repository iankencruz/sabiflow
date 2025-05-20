
-- +goose Up
ALTER TABLE auth.sessions
ADD COLUMN user_id INTEGER NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE auth.sessions
DROP COLUMN user_id;

