package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hamziCodes/CTXbank/pkg/types"
)

const (
	MemoryBankDir = "memory-bank"
	CheckpointsDir = ".state/checkpoints"
)

// Default templates conforming to the Cline Memory Bank format and token discipline.
// These are intentionally project-agnostic: they describe THIS project's structure,
// never CTXbank itself. Run 'ctx audit --apply' or `ctx prompt-sync` to fill them
// with real, project-specific content.
var defaultTemplates = map[string]string{
	"projectbrief.md": `# Project Brief

## Overview
<!-- TODO: one paragraph — what does this project build, and why? -->

## Core Deliverables
<!-- TODO: bullet list of the main things this project ships -->

## Constraints
<!-- TODO: technical or product constraints (performance, platforms, compliance) -->

## Non-Goals
<!-- TODO: what this project explicitly does NOT do -->
`,

	"productContext.md": `# Product Context

## Why This Exists
<!-- TODO: the problem this project solves, in plain language -->

## Target Audience
<!-- TODO: who uses this, and what do they need? -->

## User Experience Goals
<!-- TODO: how should using this feel? -->
`,

	"systemPatterns.md": `# System Patterns

## Core Architecture
<!-- TODO: layers, modules, and how they relate. Run 'ctx audit --apply' to seed this. -->

## Critical Design Rules
<!-- TODO: invariants that must never be broken -->
`,

	"techContext.md": `# Tech Context

## Tech Stack
<!-- TODO: languages, frameworks, key dependencies. Run 'ctx audit --apply' to seed this. -->

## Build & Run
<!-- TODO: the commands a new developer needs: build, test, lint, run -->
`,

	"decisionLog.md": `# Decision Log

All architecture and design decisions are logged here sequentially. This file is append-only.

## [INIT] Initialized CTXbank Memory Bank
- **Date:** Initial Project Bootstrap
- **Decision:** Adopted the standard Cline memory bank file layout for this project.
- **Rationale:** Human-readable files any AI coding agent already understands.
`,

	"activeContext.md": `# Active Context — updated initial_bootstrap

## Focus
<!-- TODO: what are you working on right now? One line. -->

## Recent (last 3 checkpoints)
- Initialized memory bank layout.

## Next steps
1. Run 'ctx audit --apply' to seed techContext.md and systemPatterns.md from the codebase.
2. Or run 'ctx prompt-sync' and hand the generated prompt to your AI agent.
3. Fill in projectbrief.md and productContext.md with real project intent.

## Open decisions
- None currently pending.
`,

	"progress.md": `# Progress & Milestone Ledger

## Status Overview
- Current Phase: <!-- TODO -->
- Stability: <!-- TODO -->

## Milestones
<!-- TODO: what is done, what is in progress, what is next -->
- [ ] M1: <!-- first milestone -->
`,
}

// InitBank initializes a new memory-bank directory layout in the target workspace.
// If force is false and memory-bank already exists, it returns an error to prevent clobbering.
func InitBank(workspaceDir string, force bool) error {
	bankDir := filepath.Join(workspaceDir, MemoryBankDir)
	stateDir := filepath.Join(bankDir, ".state")
	ckptDir := filepath.Join(bankDir, CheckpointsDir)

	if !force {
		if _, err := os.Stat(bankDir); err == nil {
			return fmt.Errorf("memory-bank already exists at %s; use force to reinitialize", bankDir)
		}
	}

	if err := os.MkdirAll(ckptDir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoints directory %s: %w", ckptDir, err)
	}

	// Scaffold markdown files
	for filename, content := range defaultTemplates {
		targetPath := filepath.Join(bankDir, filename)
		if !force {
			if _, err := os.Stat(targetPath); err == nil {
				continue // Preserve existing file
			}
		}
		if err := WriteAtomic(targetPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
	}

	// Create and populate initial manifest.json
	manifest := LoadOrCreateManifest(bankDir)
	initCkpt := fmt.Sprintf("ckpt_init_%s", time.Now().Format("20060102_150405"))

	for filename := range defaultTemplates {
		fullPath := filepath.Join(bankDir, filename)
		hashStr, bytesCount, err := ComputeFileHash(fullPath)
		if err != nil {
			return fmt.Errorf("failed to compute initial hash for %s: %w", filename, err)
		}
		manifest.Files[filename] = types.FileMeta{
			SHA256:         hashStr,
			Bytes:          bytesCount,
			LastCheckpoint: initCkpt,
			UpdatedAt:      time.Now().UTC(),
		}
	}

	if err := SaveManifest(bankDir, manifest); err != nil {
		return fmt.Errorf("failed to save initial manifest: %w", err)
	}

	_ = stateDir
	return nil
}

// LoadOrCreateManifest attempts to load manifest.json, or creates a new instance.
func LoadOrCreateManifest(bankDir string) *types.Manifest {
	m, err := LoadManifest(bankDir)
	if err != nil || m == nil {
		return types.NewManifest()
	}
	return m
}
