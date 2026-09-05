package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ctxbank/ctx/internal/core"
	"github.com/ctxbank/ctx/internal/mcp"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ctx-mcp: failed to determine working directory: %v\n", err)
		os.Exit(1)
	}

	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		// Auto-scaffold memory bank if not present
		_ = core.InitBank(cwd, false)
	}

	server := mcp.NewServer(bankDir, cwd)
	if err := server.Serve(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "ctx-mcp: server error: %v\n", err)
		os.Exit(1)
	}
}
