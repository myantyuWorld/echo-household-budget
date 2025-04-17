
-- +migrate Up
INSERT INTO categories (id, name) VALUES (1, '食費');
INSERT INTO categories (id, name) VALUES (2, '日用品');


-- +migrate Down
DELETE FROM categories WHERE id = 1;
DELETE FROM categories WHERE id = 2;
