-- +migrate Up
CREATE TABLE receipt_analyze_items (
  id SERIAL PRIMARY KEY,
  receipt_analyze_id INT NOT NULL,
  name TEXT NOT NULL,
  price INT NOT NULL,
  FOREIGN KEY (receipt_analyze_id) REFERENCES receipt_analyzes(id)
);

-- +migrate Down
DROP TABLE receipt_analyze_items;
