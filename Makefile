# CTXbank Makefile
# The Deterministic Memory Bank & Project Lifecycle Engine

BIN_DIR := ./bin
CTX_BIN := $(BIN_DIR)/ctx.exe
MCP_BIN := $(BIN_DIR)/ctx-mcp.exe
GO_CMD  := go

.PHONY: all build test lint audit clean

all: test build

build:
	@echo "Building static binaries with zero CGO..."
	CGO_ENABLED=0 $(GO_CMD) build -ldflags="-s -w" -o $(CTX_BIN) ./cmd/ctx
	CGO_ENABLED=0 $(GO_CMD) build -ldflags="-s -w" -o $(MCP_BIN) ./cmd/ctx-mcp
	@echo "Binaries built successfully in $(BIN_DIR)"

test:
	@echo "Running comprehensive unit & integration tests..."
	$(GO_CMD) test -v -cover ./...

lint:
	@echo "Running memory bank token & budget linter..."
	$(CTX_BIN) lint-memory

audit:
	@echo "Running 4-stage brownfield reconnaissance audit..."
	$(CTX_BIN) audit

clean:
	@echo "Cleaning binaries..."
	rm -rf $(BIN_DIR)
