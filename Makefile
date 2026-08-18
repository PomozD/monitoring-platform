.PHONY: help up down restart ps logs config

COMPOSE := docker compose --env-file .env -f deploy/docker/docker-compose.yml

help:
	@echo "Available commands:"
	@echo "  make up       Start infrastructure"
	@echo "  make down     Stop infrastructure"
	@echo "  make restart  Restart infrastructure"
	@echo "  make ps       Show containers"
	@echo "  make logs     Show container logs"
	@echo "  make config   Validate Docker Compose configuration"

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
