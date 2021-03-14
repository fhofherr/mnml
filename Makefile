GO := go
BIN_DIR := bin
GO_FILES := $(shell find . -iname '*.go')

.DEFAULT_GOAL := $(BIN_DIR)/mnml

.PHONY: generate
generate: $(BIN_DIR)/stringer
	PATH=$(BIN_DIR):$(PATH) $(GO) generate ./...

$(BIN_DIR)/stringer:
	GOBIN=$(abspath ./$(BIN_DIR)) $(GO) install golang.org/x/tools/cmd/stringer

$(BIN_DIR)/mnml: $(GO_FILES)
	$(GO) build -o $@ ./cmd/mnml

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
