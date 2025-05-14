-- +migrate Up
CREATE TABLE refresh_token (
  id          UUID PRIMARY KEY,
  user_id     BIGINT NOT NULL UNIQUE,
  issued_at   TIMESTAMP NOT NULL,
  expires_at  TIMESTAMP NOT NULL,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS refresh_token;
