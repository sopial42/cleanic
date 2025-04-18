-- +migrate Down

-- +migrate Up

CREATE TABLE patient (
  id        BIGSERIAL PRIMARY KEY,
  firstname TEXT      NOT NULL,
  lastname  TEXT      NOT NULL,
  email     TEXT      NOT NULL UNIQUE
);

ALTER SEQUENCE patient_id_seq RESTART WITH 10001;
