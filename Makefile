.PHONY: help build run dev lint format docs-generate migrate-up migrate-down migrate-force docker-up docker-down

help:
	@echo "Available commands"
	@echo " make build				- Build the application"
	@echo " make run 				- Run the application"
	@echo " make stop 				- Stop the application"
	@echo " make dev 				- Run the application in development mode"
	@echo " make lint 				- Run linter on the codebase"
	@echo " make format 			- Format the code and re-arrange the imports"
	@echo " make docs-generate 		- Generate swagger documentation"
	@echo " make migrate-up			- Apply database migrations"
	@echo " make migrate-down		- Rollback database migrations"
	@echo " make migrate-force		- Resets migration state to before any migration"
	@echo " make docker-up			- Run docker-compose to setup database & local cloud"
	@echo " make docker-down		- Close docker-compose"

build:
	@echo "Building all binaries"
	@mkdir -p bin
	@for cmd in cmd/*/; do \
			if [ -d "$$cmd" ]; then \
				binary=$$(basename $$cmd); \
				echo "Building $$binary...."; \
				go build -o bin/$$binary ./$$cmd; \
			fi \
		done

run:
	go run ./cmd/api

stop:
	pkill -f "./tmp/main"

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

format:
	@gofmt -s -w .

docs-generate:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

migrate-force:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" force 1

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down