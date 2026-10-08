package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hamziCodes/CTXbank/internal/core"
	"github.com/hamziCodes/CTXbank/internal/git"
)

// doctorCheck is a single self-diagnostic check result.
type doctorCheck struct {
	name   string
	status string // "ok", "warn", "fail"
	detail string
}

func runDoctor(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var checks []doctorCheck
	fail := false

	report := func(name, status, detail string) {
		checks = append(checks, doctorCheck{name, status, detail})
		if status == "fail" {
			fail = true
		}
	}

	// 1. Binary + version
	report("ctx binary", "ok", fmt.Sprintf("running version %s", Version))

	// 2. git present
	if out, err := exec.Command("git", "--version").Output(); err != nil {
		report("git in PATH", "fail", "git not found — checkpoints and status need git 2.30+")
	} else {
		report("git in PATH", "ok", string(out[:len(out)-1]))
	}

	// 3. memory-bank initialized
	bankDir := filepath.Join(cwd, core.MemoryBankDir)
	if _, err := os.Stat(bankDir); os.IsNotExist(err) {
		report("memory-bank", "fail", "not initialized here — run 'ctx init' first")
	} else {
		report("memory-bank", "ok", bankDir)

		// 4. manifest loads and hashes match (no external drift)
		manifest, err := core.LoadManifest(bankDir)
		if err != nil || manifest == nil {
			report("manifest", "fail", fmt.Sprintf("cannot load manifest.json: %v", err))
		} else {
			drifted := []string{}
			for name, meta := range manifest.Files {
				full := filepath.Join(bankDir, name)
				hash, _, err := core.ComputeFileHash(full)
				if err != nil {
					drifted = append(drifted, name+" (missing)")
					continue
				}
				if hash != meta.SHA256 {
					drifted = append(drifted, name)
				}
			}
			if len(drifted) > 0 {
				report("manifest drift", "warn", fmt.Sprintf("files changed outside ctx: %v — run 'ctx pause' to re-sync", drifted))
			} else {
				report("manifest drift", "ok", fmt.Sprintf("%d files match manifest hashes", len(manifest.Files)))
			}
		}
	}

	// 5. repo safety (rebase / detached / conflicts)
	if err := git.ValidateRepoSafety(cwd); err != nil {
		report("repo state", "warn", err.Error())
	} else {
		report("repo state", "ok", "no rebase, detached HEAD, or merge conflicts")
	}

	// Render
	fmt.Println("CTXbank doctor — environment self-check")
	fmt.Println()
	for _, c := range checks {
		mark := "✓"
		switch c.status {
		case "warn":
			mark = "!"
		case "fail":
			mark = "✗"
		}
		fmt.Printf("  %s %-14s %s\n", mark, c.name, c.detail)
	}
	fmt.Println()
	if fail {
		fmt.Println("Result: issues found — fix the ✗ items above.")
		os.Exit(1)
	}
	fmt.Println("Result: all critical checks passed.")
	fmt.Println()
	fmt.Println("MCP snippet for this project (Cursor / Claude Code / Antigravity):")
	fmt.Printf(`{
  "mcpServers": {
    "ctxbank": {
      "command": "ctx",
      "args": ["serve", "--mcp"],
      "cwd": %q
    }
  }
}
`, cwd)
}
