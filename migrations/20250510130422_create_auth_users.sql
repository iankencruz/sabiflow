-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
  id SERIAL PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

);

-- +goose Down
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth;
