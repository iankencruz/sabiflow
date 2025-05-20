
-- +goose Up
CREATE TABLE auth.session_data (
  token TEXT NOT NULL REFERENCES auth.sessions(token) ON DELETE CASCADE,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  PRIMARY KEY (token, key)
);

-- +goose Down
DROP TABLE IF EXISTS auth.session_data;

