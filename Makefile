include MANIFEST

export

PROJECT_DIR = $(shell pwd)

TOOLS_DIR = $(PROJECT_DIR)/tools
TOOLS_BIN_DIR = $(TOOLS_DIR)/bin
$(shell [ -f $(TOOLS_BIN_DIR) ] || mkdir -p $(TOOLS_BIN_DIR))

GO ?= go

.PHONY: run-app

run-app:
	$(GO) run ./cmd/app

GO_LINT_TOOL = $(TOOLS_BIN_DIR)/golangci-lint

.PHONY: .install-golangci-lint
.install-golangci-lint:
	@[ -f $(GO_LINT_TOOL) ] \
	|| GOBIN=$(TOOLS_BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION)

.PHONY: lint
lint: .install-golangci-lint
	$(GO_LINT_TOOL) run ./... --fix --config=./.golangci.yml

GO_BUILD_PATH = $(PROJECT_DIR)/bin
GO_BUILD_APP_PATH = $(GO_BUILD_PATH)/app

.PHONY: build
build:
	$(GO) build -o $(GO_BUILD_APP_PATH) ./cmd/app

.PHONY: run
run:
	$(GO) run ./cmd/app