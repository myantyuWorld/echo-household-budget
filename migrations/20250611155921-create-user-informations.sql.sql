-- +migrate Up
CREATE TABLE IF NOT EXISTS user_informations (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  information_id INTEGER NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (user_id, information_id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_informations;
