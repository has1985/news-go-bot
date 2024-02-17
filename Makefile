PROJECT_ROOT            := github.com/glocksoftware/news-go-bot
SRCROOT_ON_HOST         := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
SRCROOT_IN_CONTAINER    := /go/src/$(PROJECT_ROOT)

# configuration for the database postgres
WITH_DATABASE           = true
DATABASE_ADDRESS        ?= localhost:5435
DATABASE_USERNAME       ?= postgres
DATABASE_PASSWORD       ?= postgres
DATABASE_NAME           = news-go-bot
DATABASE_URL            ?= postgres://$(DATABASE_USERNAME):$(DATABASE_PASSWORD)@$(DATABASE_ADDRESS)/$(DATABASE_NAME)?sslmode=disable

MIGRATETOOL_IMAGE           = infoblox/migrate:latest
MIGRATION_PATH_IN_CONTAINER = $(SRCROOT_IN_CONTAINER)/migrations

DOCKER_RUNNER           := docker run -u `id -u`:`id -g` --rm
DOCKER_RUNNER           += -v $(SRCROOT_ON_HOST):$(SRCROOT_IN_CONTAINER)

.PHONY migrate-up: migrate-up-atlas
migrate-up-atlas:
	@$(DOCKER_RUNNER) --net="host" $(MIGRATETOOL_IMAGE) --verbose --path=$(MIGRATION_PATH_IN_CONTAINER)/ --database.dsn=$(DATABASE_URL) up $(UP)


.PHONY migrate-down: migrate-down-atlas
migrate-down-atlas:
	@$(DOCKER_RUNNER) --net="host" $(MIGRATETOOL_IMAGE) --verbose --path=$(MIGRATION_PATH_IN_CONTAINER)/ --database.dsn=$(DATABASE_URL) down $(DOWN)


UP ?= 1
DOWN ?= 1

.PHONY env-up:
env-up:
	docker-compose up --remove-orphans -d

.PHONY env-down:
env-down:
	docker-compose down

.PHONY env-stop:
env-stop:
	docker-compose stop