GO := go
PRE_COMMIT := pre-commit

BIN_DIR := bin
GO_FILES := $(shell find . -iname '*.go')

.DEFAULT_GOAL := $(BIN_DIR)/mnml

.PHONY: generate
generate: $(BIN_DIR)/stringer ## Run Go generators.
	PATH=$(BIN_DIR):$(PATH) $(GO) generate ./...

.PHONY: lint
lint: ## Run linter on all files.
	$(PRE_COMMIT) run --all-files

.PHONY: release
release: $(BIN_DIR)/goreleaser ## Create a release using goreleaser.
	$< release --rm-dist

.PHONY: clean
clean: ## Clean directory.
	rm -rf $(BIN_DIR)

.PHONY: help
help: ## Show this help.
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

$(BIN_DIR)/stringer:
	GOBIN=$(abspath ./$(BIN_DIR)) $(GO) install golang.org/x/tools/cmd/stringer

$(BIN_DIR)/goreleaser:
	GOBIN=$(abspath ./$(BIN_DIR)) $(GO) install github.com/goreleaser/goreleaser

$(BIN_DIR)/mnml: $(GO_FILES)
	$(GO) build -o $@ ./cmd/mnml
