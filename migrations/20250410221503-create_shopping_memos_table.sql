-- +migrate Up
CREATE TABLE IF NOT EXISTS shopping_memos (
    id SERIAL PRIMARY KEY,
    household_book_id INTEGER NOT NULL REFERENCES household_books(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL DEFAULT 0,
    memo TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shopping_memos_household_book_id ON shopping_memos(household_book_id);
CREATE INDEX idx_shopping_memos_category_id ON shopping_memos(category_id);
CREATE INDEX idx_shopping_memos_is_completed ON shopping_memos(is_completed);

-- +migrate Down
DROP TABLE IF EXISTS shopping_memos; 