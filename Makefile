# ビルド
up-build:
	docker compose up -d --build
# 起動
up:
	docker compose up -d
# 停止
down:
	docker compose down
# 再起動
restart:
	@make down
	@make up

# マイグレーション
dev-stat:
	sql-migrate status -env="development"

dev-up:
	sql-migrate up -env="development"

dev-down:
	sql-migrate down -env="development"

exec-db:
	docker compose exec db /bin/bash

.PHONY: run test migrate-up migrate-down

run:
	go run cmd/main.go

test:
	go test -v ./...

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down
