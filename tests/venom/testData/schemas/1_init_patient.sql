-- +migrate Down

-- +migrate Up

CREATE TABLE patient(
  id        INT GENERATED ALWAYS AS IDENTITY (START WITH 10001) PRIMARY KEY,
  firstname TEXT NOT NULL,
  lastname  TEXT NOT NULL,  
  email     TEXT NOT NULL UNIQUE
);

