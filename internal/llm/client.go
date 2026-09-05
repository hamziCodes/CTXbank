package llm

import (
	"errors"
)

// TargetDoc defines the destination memory bank file for an ingested chunk.
type TargetDoc string

const (
	DocProjectBrief   TargetDoc = "projectbrief.md"
	DocProductContext TargetDoc = "productContext.md"
	DocSystemPatterns TargetDoc = "systemPatterns.md"
	DocTechContext    TargetDoc = "techContext.md"
	DocDiscard        TargetDoc = "discard"
)

var (
	ErrLLMUnavailable = errors.New("local LLM (Ollama) is offline or unreachable; falling back to deterministic Tier-0")
)

// Client defines the interface for local Tier-1 intelligence.
type Client interface {
	IsAvailable() bool
	ClassifyChunk(text string) (TargetDoc, error)
	SummarizeNarrative(contextData string) (string, error)
}
