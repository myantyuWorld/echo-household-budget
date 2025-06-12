-- +migrate Up
CREATE TYPE information_category AS ENUM ('bug_report', 'feature_request', 'other');

CREATE TABLE IF NOT EXISTS informations (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  category information_category NOT NULL,
  is_published BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS informations;

DROP TYPE IF EXISTS information_category;
