
-- +migrate Up
ALTER TABLE user_accounts ADD COLUMN picture_url VARCHAR(255);

-- +migrate Down
ALTER TABLE user_accounts DROP COLUMN picture_url;
