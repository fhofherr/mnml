GO := go

PROJECT_ROOT:= $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BIN_DIR := $(PROJECT_ROOT)bin

.PHONY: generate
generate: $(BIN_DIR)/stringer
	PATH=$(BIN_DIR):$(PATH) $(GO) generate ./...

$(BIN_DIR)/stringer:
	GOBIN=$(BIN_DIR) $(GO) install golang.org/x/tools/cmd/stringer
