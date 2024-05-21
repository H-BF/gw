GO_BIN ?= $(shell which go)
GOLANGCI_BIN ?= $(shell which golangci-lint)
DOCKER_BIN ?= $(shell which docker)
MIGRATE_BIN ?= $(shell which migrate)

APP_DIR ?= cmd
APP_NAME := gw
APP_PATH := $(APP_DIR)/$(APP_NAME).go

.PHONY: help
help: ## command tips
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: go-deps
go-deps: ## install dependencies from go.mod
	@echo install go dependecies... && \
	$(GO_BIN) mod tidy && \
	$(GO_BIN) mod vendor && \
	$(GO_BIN) mod verify && \
	echo all dependecies installed

.PHONY: run
run: ## run app with go runtime
	@echo run app on dev mode && \
	$(GO_BIN) run $(APP_PATH) && \
	echo app stopped

PG_USER ?=
PG_PWD ?=
PG_PORT ?=
.PHONY: up-pg
up-pg: ## up postgres docker image for casbin adapter
	@echo up postgres in docker && \
	$(DOCKER_BIN) run \
		--name gw-casbin-adapter  \
		-e POSTGRES_USER=$(PG_USER) \
		-e POSTGRES_PASSWORD=${PG_PWD} \
		-e POSTGRES_DB=casbin \
		-p $(PG_PORT):5432 \
		-d postgres && \
	ehco postgres listen on 127.0.0.1:$(PG_PORT)

.PHONY: .install-goose
.install-goose: ## install goose migrate tool
ifeq ($(filter $(GOOSE_CUR_VERSION), $(GOOSE_LATEST_VERSION)),)
	@echo installing \'goose\' util... && \
	$(GO_BIN) install github.com/pressly/goose/v3/cmd/goose@latest
else
	@echo >/dev/null
endif

PG_MIGRATIONS?=$(CURDIR)/internal/registry/casbin/pg/migrations
PG_URI?=
.PHONY: casbin-migrate-up
casbin-migrate-up: ## run Casbin Postgres migrations
ifneq ($(PG_URI),)
	@$(MAKE) .install-goose && \
	cd $(PG_MIGRATIONS) && \
	goose -table=casbin postgres $(PG_URI) up
else
	$(error need define PG_URI environment variable)
endif	

.PHONY: lint
lint: ## run full lint
	@echo full lint... && \
	$(GOLANGCI_BIN) cache clean && \
	$(GOLANGCI_BIN) run \
		--timeout=120s \
		--config=$(CURDIR)/.golangci.yaml \
		-v $(CURDIR)/... && \
	echo -=OK=-