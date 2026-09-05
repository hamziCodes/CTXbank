package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ctxbank/ctx/pkg/types"
)

const (
	MemoryBankDir = "memory-bank"
	CheckpointsDir = ".state/checkpoints"
)

// Default templates conforming to the Cline Memory Bank format and token discipline
var defaultTemplates = map[string]string{
	"projectbrief.md": `# Project Brief

## Overview
High-level description of the project vision, core deliverables, and primary architectural goals.

## Core Constraints
- Single static binary footprint (< 15MB) with zero runtime dependencies.
- Atomic crash-safe writes on all file mutations.
- Strict token discipline: keep hot context files condensed (< 150 lines).

## Non-Goals
- Cloud-hosted state storage (local-first design).
- Unconfirmed automatic file overwrites.
`,

	"productContext.md": `# Product Context

## Why This Exists
Provides a deterministic, crash-safe, local-first memory layer for AI coding agents and developers.
Replaces brittle prompt-engineering conventions with cryptographic verification, atomic file operations, and MCP delta-sync.

## Target Audience
- Developers collaborating with AI coding agents (Cline, Cursor, Claude Code, Antigravity).
- Engineers needing zero-drift pickup briefs across multi-day context switches.

## User Experience Goals
- Zero friction: commands work instantaneously without requiring cloud credentials.
- Zero token waste: returns 20-byte {"unchanged": true} when context has not drifted.
`,

	"systemPatterns.md": `# System Patterns

## Core Architecture
- **Layer 1: CLI (` + "`cmd/ctx`" + `)**: Static binary for terminal workflows.
- **Layer 2: Core Storage (` + "`internal/core`" + `)**: Atomic writes, SHA-256 manifest engine.
- **Layer 3: Git Porcelain (` + "`internal/git`" + `)**: Non-destructive status and diff analysis.
- **Layer 4: MCP Protocol (` + "`internal/mcp`" + `)**: Stdio-based agent protocol with session delta caching.

## Critical Design Rules
- Atomic write pattern: write to ` + "`<file>.tmp.<rand>`" + ` -> ` + "`fsync()`" + ` -> ` + "`rename()`" + `.
- Never store API keys in plaintext config files (resolve strictly via environment variables).
`,

	"techContext.md": `# Tech Context

## Tech Stack
- **Language:** Go 1.22+
- **Build Target:** Statically linked binary (` + "`CGO_ENABLED=0`" + `)
- **Git Integration:** Porcelain CLI queries (` + "`git status --porcelain`" + `, ` + "`git diff`" + `)
- **Agent Protocol:** Model Context Protocol (MCP) over Stdio
- **Local LLM Tier:** Ollama HTTP REST API (` + "`localhost:11434`" + `)

## Runtime Requirements
- OS: Windows, Linux, macOS
- Git: 2.30+ installed in PATH
`,

	"decisionLog.md": `# Decision Log

All architecture and design decisions are logged here sequentially. This file is append-only.

## [INIT] Initialized CTXbank Memory Bank
- **Date:** Initial Project Bootstrap
- **Decision:** Adopt standard Cline memory bank naming schema with deterministic Go engine underneath.
- **Rationale:** Interoperates with established ecosystem conventions while eliminating amnesia and file corruption.
`,

	"activeContext.md": `# Active Context — updated initial_bootstrap

## Focus
System initialization and baseline scaffolding.

## Recent (last 3 checkpoints)
- Initialized CTXbank deterministic memory bank layout.

## Next steps
1. Complete core storage and manifest verification tests.
2. Build and verify Git porcelain integration and safe checkpointing engine.
3. Validate CLI commands (` + "`init`" + `, ` + "`status`" + `, ` + "`pause`" + `, ` + "`resume`" + `).

## Open decisions
- None currently pending.
`,

	"progress.md": `# Progress & Milestone Ledger

## Status Overview
- Current Phase: Milestone 1 — Core Foundation & Deterministic CLI
- Stability: Initial Development

## Milestones
- [ ] M1: Core Foundation & Deterministic CLI
  - [x] Phase 01: Core Manifest & Atomic Store
  - [ ] Phase 02: Git Porcelain Integration & Safe Checkpoint Engine
  - [ ] Phase 03: Primary CLI Surface (init, status, pause, resume)
- [ ] M2: MCP Server & Multi-Agent Interop
- [ ] M3: Reconnaissance Audit & CI Budget Linter
- [ ] M4: Research Ingestion Pipeline & Local LLM Integration
- [ ] M5: Cross-Project Workspace Dashboard & Hardening
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
