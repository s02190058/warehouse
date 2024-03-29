include .env.local

include MANIFEST

export

PROJECT_DIR = $(shell pwd)

GO ?= go

# tools

TOOLS_DIR = $(PROJECT_DIR)/tools
TOOLS_BIN_DIR = $(TOOLS_DIR)/bin
$(shell [ -f $(TOOLS_BIN_DIR) ] || mkdir -p $(TOOLS_BIN_DIR))

GO_LINT_TOOL = $(TOOLS_BIN_DIR)/golangci-lint
.PHONY: .install-golangci-lint
.install-golangci-lint:
	@[ -f $(GO_LINT_TOOL) ] \
	|| GOBIN=$(TOOLS_BIN_DIR) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: .install-golangci-lint
	$(GO_LINT_TOOL) run ./... --fix --config=./.golangci.yml

GO_SWAG_TOOL = $(TOOLS_BIN_DIR)/swag

.PHONY: .install-swag
.install-swag:
	@[ -f $(GO_SWAG_TOOL) ] \
	|| GOBIN=$(TOOLS_BIN_DIR) $(GO) install github.com/swaggo/swag/cmd/swag@v$(SWAG_VERSION)

.PHONY: swag
swag: .install-swag
	$(GO_SWAG_TOOL) fmt \
	&& $(GO_SWAG_TOOL) init -g ./internal/transport/http/router.go

GO_MIGRATE_TOOL = $(TOOLS_BIN_DIR)/migrate

.PHONY: install-migrate
install-migrate:
	@[ -f $(GO_MIGRATE_TOOL) ] \
	|| GOBIN=$(TOOLS_BIN_DIR) $(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v$(MIGRATE_VERSION)

POSTGRES_URL='postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable'

.PHONY: migrate-up
migrate-up:
	$(GO_MIGRATE_TOOL) -database $(POSTGRES_URL) -path migrations/warehouse up

.PHONY: migrate-down
migrate-down:
	$(GO_MIGRATE_TOOL) -database $(POSTGRES_URL) -path migrations/warehouse down

# local build

GO_BUILD_PATH = $(PROJECT_DIR)/bin
GO_BUILD_WAREHOUSE_PATH = $(GO_BUILD_PATH)/warehouse

.PHONY: build
build:
	CGO_ENABLED=0 $(GO) build -o $(GO_BUILD_WAREHOUSE_PATH) ./cmd/warehouse

.PHONY: run
run:
	CGO_ENABLED=0 $(GO) run ./cmd/warehouse

.PHONY: test
test:
	CGO_ENABLED=1 $(GO) test -race -v ./...

.PHONY: test-all
test-all:
	CGO_ENABLED=1 $(GO) test -race -tags integration -v ./...

# docker build

DOCKER ?= docker

.PHONY: up
up:
	$(DOCKER) compose up -d --build

.PHONY: down
down:
	$(DOCKER) compose down

.PHONY: rm-pg-data
rm-pg-data:
	$(DOCKER) volume rm warehouse_pg-data