-- +goose Up
CREATE TABLE auth.sessions (
  token TEXT PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS auth.sessions;
