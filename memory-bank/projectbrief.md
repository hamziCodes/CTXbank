# Project Brief

## Overview
High-level description of the project vision, core deliverables, and primary architectural goals.

## Core Constraints
- Single static binary footprint (< 15MB) with zero runtime dependencies.
- Atomic crash-safe writes on all file mutations.
- Strict token discipline: keep hot context files condensed (< 150 lines).

## Non-Goals
- Cloud-hosted state storage (local-first design).
- Unconfirmed automatic file overwrites.
