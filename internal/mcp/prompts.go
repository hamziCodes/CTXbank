package mcp

// PromptDescriptor defines an MCP prompt template.
type PromptDescriptor struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListPrompts returns available MCP prompts.
func ListPrompts() []PromptDescriptor {
	return []PromptDescriptor{
		{
			Name:        "resume-session",
			Description: "Pre-loads active context and progress into agent first turn, generating a pickup brief.",
		},
	}
}

// GetResumePrompt generates the initial prompt message for an agent session.
func GetResumePrompt(activeCtx, progress string) string {
	return `# Agent Pickup Brief & Ground Truth Context

Before modifying code, review the current focus and progress:

---
## Current Active Context
` + activeCtx + `

---
## Milestone Progress
` + progress + `

Please adhere strictly to the rules in AGENTS.md, maintain atomic file writes, and update activeContext.md upon completion.
`
}
