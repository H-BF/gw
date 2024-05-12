GO_BIN ?= $(shell which go)
GOLANGCI_BIN ?= $(shell which golangci-lint)

APP_DIR ?= cmd
APP_NAME := sgroups-k8s-adapter
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

.PHONY: lint
lint: ## run full lint
	@echo full lint... && \
	$(GOLANGCI_BIN) cache clean && \
	$(GOLANGCI_BIN) run \
		--timeout=120s \
		--config=$(CURDIR)/.golangci.yaml \
		-v $(CURDIR)/... && \
	echo -=OK=-