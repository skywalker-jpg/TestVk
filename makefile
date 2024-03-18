GO_BUILD_APP_PATH := bin

DOCKER_COMPOSE_FILE := docker-compose.yml

.PHONY: up
up:
	@echo "Starting Docker containers..."
	docker-compose up -d --build

.PHONY: down
down:
	@echo "Stopping Docker containers..."
	docker compose down

