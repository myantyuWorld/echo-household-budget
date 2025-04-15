
-- +migrate Up
CREATE TABLE IF NOT EXISTS shopping_amounts (
    id SERIAL PRIMARY KEY,
    household_book_id INTEGER NOT NULL REFERENCES household_books(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL DEFAULT 0,
    date DATE NOT NULL,
    memo TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_book_id) REFERENCES household_books(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_shopping_amounts_household_book_id ON shopping_amounts(household_book_id);
CREATE INDEX idx_shopping_amounts_category_id ON shopping_amounts(category_id);

-- +migrate Down
DROP TABLE IF EXISTS shopping_amounts;
