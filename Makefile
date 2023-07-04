-include .env

.PHONY: up
up:
	$(eval SERVICE = ${s})
	@docker-compose up -d --no-build --remove-orphans ${SERVICE}
	@docker-compose ps

down:
	@docker-compose down --remove-orphans --volumes

migrate:
	@go run cmd/setup_db/main.go