DC = docker-compose
IMAGE_NAME := soliant-mock-api
VERSION := latest

up:
	$(DC) up -d

up-db-local:
	$(DC) -f local-docker-compose.yml up -d

up-db:
	$(DC) up -d db

build:
	$(DC) build

stop:
	$(DC) stop

down:
	$(DC) down

truncate:
	$(DC) -f local-docker-compose.yml down -v

logs:
	$(DC) logs

db-shell:
	$(DC) exec db psql -U soliant-mock-api

create-volume:
	docker volume create --name=postgres_data

remove-volume:
	docker volume rm postgres_data

reestart: stop down build up logs

.PHONY: up build stop down logs db-shell create-volume remove-volume
