
-- +migrate Up
CREATE TABLE IF NOT EXISTS category_limits (
    id SERIAL PRIMARY KEY,
    household_book_id INTEGER NOT NULL REFERENCES household_books(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    limit_amount INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  
    FOREIGN KEY (household_book_id) REFERENCES household_books(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_category_limits_household_book_id ON category_limits(household_book_id);
CREATE INDEX idx_category_limits_category_id ON category_limits(category_id);

-- +migrate Down
DROP TABLE IF EXISTS category_limits;