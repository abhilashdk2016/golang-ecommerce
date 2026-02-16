DB_URL=postgresql://postgres:password@localhost:5433/golang_ecommerce?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

build:
	@echo "Building all binaries...."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
    		if [ -d "$$cmd" ]; then \
    			binary=$$(basename $$cmd); \
    			echo "Building $$binary..."; \
    			go build -o bin/$$binary ./$$cmd; \
    		fi \
    	done

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

docs-generate:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

docker-up:
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	docker-compose -f docker/docker-compose.yml down

.PHONY: build run dev lint migrate-up migrate-down docker-up docker-down format docs-generate