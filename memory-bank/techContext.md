# Tech Context

## Tech Stack
- **Language:** Go 1.22+
- **Build Target:** Statically linked binary (`CGO_ENABLED=0`)
- **Git Integration:** Porcelain CLI queries (`git status --porcelain`, `git diff`)
- **Agent Protocol:** Model Context Protocol (MCP) over Stdio
- **Local LLM Tier:** Ollama HTTP REST API (`localhost:11434`)

## Runtime Requirements
- OS: Windows, Linux, macOS
- Git: 2.30+ installed in PATH
