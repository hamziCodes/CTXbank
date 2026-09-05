package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctxbank/ctx/internal/core"
)

// Config holds runtime settings for CTXbank operations.
type Config struct {
	OllamaURL    string   `json:"ollama_url"`
	DefaultModel string   `json:"default_model"`
	ScanPaths    []string `json:"scan_paths"`
}

const (
	ConfigDirName  = ".context"
	ConfigFileName = "config.toml"
)

// DefaultConfig returns baseline configuration values.
func DefaultConfig() *Config {
	ollamaURL := os.Getenv("OLLAMA_HOST")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2"
	}

	return &Config{
		OllamaURL:    ollamaURL,
		DefaultModel: model,
		ScanPaths:    []string{"."},
	}
}

// LoadConfig reads .context/config.toml from repoDir if present, resolving env vars.
func LoadConfig(repoDir string) (*Config, error) {
	cfg := DefaultConfig()
	configPath := filepath.Join(repoDir, ConfigDirName, ConfigFileName)

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		switch key {
		case "ollama_url":
			cfg.OllamaURL = val
		case "default_model":
			cfg.DefaultModel = val
		}
	}

	// Environment variable overrides (BYOK guarantee: env takes highest precedence)
	if envURL := os.Getenv("OLLAMA_HOST"); envURL != "" {
		cfg.OllamaURL = envURL
	}
	if envModel := os.Getenv("OLLAMA_MODEL"); envModel != "" {
		cfg.DefaultModel = envModel
	}

	return cfg, nil
}

// SaveDefaultConfig creates a clean .context/config.toml template with zero plaintext keys.
func SaveDefaultConfig(repoDir string) error {
	configDir := filepath.Join(repoDir, ConfigDirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return nil // Do not overwrite existing user config
	}

	content := `# CTXbank Local Configuration (.context/config.toml)
# Note: Sensitive keys are NEVER written here.
# Configure them via environment variables:
#   export OLLAMA_HOST="http://localhost:11434"
#   export OLLAMA_MODEL="llama3.2"

ollama_url = "http://localhost:11434"
default_model = "llama3.2"
`
	return core.WriteAtomic(configPath, []byte(content), 0644)
}
