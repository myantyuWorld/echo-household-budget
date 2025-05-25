-- +migrate Up
CREATE TYPE analyze_status AS ENUM ('pending', 'finished');

CREATE TABLE receipt_analyzes (
  id SERIAL PRIMARY KEY,
  image_url TEXT NOT NULL,
  analyze_status analyze_status NOT NULL,
  total_price INT NOT NULL,
  household_book_id INT NOT NULL,
  FOREIGN KEY (household_book_id) REFERENCES household_books(id)
);

-- +migrate Down
DROP TABLE receipt_analyzes;

DROP TYPE analyze_status;
