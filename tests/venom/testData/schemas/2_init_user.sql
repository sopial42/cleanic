-- +migrate Down

-- +migrate Up

CREATE TABLE users (
  id        BIGSERIAL PRIMARY KEY,
  password  TEXT      NOT NULL CHECK (password <> ''),
  email     TEXT      NOT NULL UNIQUE,
  roles     JSONB NOT NULL
);

ALTER SEQUENCE users_id_seq RESTART WITH 10001;
