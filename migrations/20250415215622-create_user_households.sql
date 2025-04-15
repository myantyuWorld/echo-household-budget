
-- +migrate Up
CREATE TABLE IF NOT EXISTS user_households (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    household_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (household_id) REFERENCES household_books(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_households_user_id ON user_households(user_id);
CREATE INDEX idx_user_households_household_id ON user_households(household_id);

-- +migrate Down
DROP TABLE IF EXISTS user_households;