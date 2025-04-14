-- +migrate Up
CREATE TABLE IF NOT EXISTS household_books (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_household_books_user_id ON household_books(user_id);
CREATE INDEX idx_household_books_created_at ON household_books(created_at);

-- +migrate Down
DROP TABLE IF EXISTS household_books; 