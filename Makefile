GO ?= go
PRE_COMMIT ?= pre-commit

TOOLS_GO := tools.go
TOOL_PKGS := $(shell awk '/[[:space:]]+_[[:space:]]+"[^"]+"/ { gsub(/"/, "", $$2); print $$2}' $(TOOLS_GO))
TOOLS_DIR := .tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin
TOOL_BINARIES := $(patsubst %,$(TOOLS_BIN_DIR)/%,$(notdir $(TOOL_PKGS)))

CMD_DIR := cmd
BIN_DIR := bin

GO_FILES := $(shell find . -iname '*.go')
MAIN_GO_FILES := $(wildcard $(CMD_DIR)/*/main.go)
CMD_PKGS := $(patsubst $(CMD_DIR)/%/main.go,$(CMD_DIR)/%,$(MAIN_GO_FILES))
CMD_BINARIES := $(patsubst $(CMD_DIR)/%/main.go,$(BIN_DIR)/%,$(MAIN_GO_FILES))

.DEFAULT_GOAL := build

.PHONY: build
build: $(CMD_BINARIES) ## Build all command binaries

.PHONY: generate
generate: $(TOOL_BINARIES) $(GO_FILES) ## Run Go generators.
	PATH=$(abspath ./$(TOOLS_BIN_DIR)):$(PATH) $(GO) generate ./...

.PHONY: lint
lint: $(GO_FILES) ## Run linter on all files.
	$(PRE_COMMIT) run --all-files

.PHONY: test
test: ## Run all tests
	$(GO) test ./...

.PHONY: install
install: $(CMD_PKGS) ## Install the binary to the default GOBIN
	$(foreach pkg,$^,$(GO) install ./$(pkg))

.PHONY: release
release: $(TOOLS_BIN_DIR)/goreleaser ## Create a release using goreleaser.
	$< release --rm-dist

.PHONY: clean
clean: ## Clean directory.
	rm -rf $(BIN_DIR)
	rm -rf $(TOOLS_DIR)

.PHONY: help
help: ## Show this help.
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

$(BIN_DIR)/%: $(CMD_DIR)/%/main.go $(GO_FILES)
	$(GO) build -o $@ ./$(dir $<)

$(TOOLS_BIN_DIR)/%: $(TOOLS_GO)
	GOBIN=$(abspath ./$(TOOLS_BIN_DIR)) $(GO) install $(shell grep '_ $%' $(TOOLS_GO) | cut -d'"' -f2)
