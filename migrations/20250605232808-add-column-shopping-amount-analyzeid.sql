-- +migrate Up
alter table
  shopping_amounts
add
  column analyze_id int not null default 0;

-- +migrate Down
alter table
  shopping_amounts drop column analyze_id;
