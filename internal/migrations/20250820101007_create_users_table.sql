-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX users_email_unique ON users (lower(email));

-- +goose Down
DROP INDEX IF EXISTS users_email_unique;
DROP TABLE IF EXISTS users;