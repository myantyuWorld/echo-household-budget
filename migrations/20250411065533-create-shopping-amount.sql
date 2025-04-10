
-- +migrate Up
CREATE TABLE IF NOT EXISTS shopping_amounts (
    id SERIAL PRIMARY KEY,
    household_book_id INTEGER NOT NULL REFERENCES household_books(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL DEFAULT 0,
    date DATE NOT NULL,
    memo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS shopping_amounts;
