ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: help up down restart ps logs config migrate-up migrate-down migrate-version

COMPOSE := docker compose --env-file .env -f deploy/docker/docker-compose.yml

MIGRATIONS_PATH := services/auth-service/migrations

DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make up               Start infrastructure"
	@echo "  make down             Stop infrastructure"
	@echo "  make restart          Restart infrastructure"
	@echo "  make ps               Show containers"
	@echo "  make logs             Show container logs"
	@echo "  make config           Validate Docker Compose configuration"
	@echo "  make migrate-up       Apply database migrations"
	@echo "  make migrate-down     Rollback last migration"
	@echo "  make migrate-version  Show migration version"

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) restart

ps:
	$(COMPOSE) ps

logs:
	$(COMPOSE) logs -f

config:
	$(COMPOSE) config

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version