-include .env

.PHONY: up
up:
	$(eval SERVICE = ${s})
	@docker-compose up -d --no-build --remove-orphans ${SERVICE}
	@docker-compose ps
