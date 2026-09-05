package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigDefault(t *testing.T) {
	tempDir := t.TempDir()

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.OllamaURL == "" || cfg.DefaultModel == "" {
		t.Errorf("expected defaults populated, got: %+v", cfg)
	}
}

func TestLoadConfigFileAndEnvOverride(t *testing.T) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ConfigDirName)
	_ = os.MkdirAll(configDir, 0755)

	tomlContent := `ollama_url = "http://custom-host:8080"` + "\n" + `default_model = "mistral"` + "\n"
	_ = os.WriteFile(filepath.Join(configDir, ConfigFileName), []byte(tomlContent), 0644)

	cfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.OllamaURL != "http://custom-host:8080" {
		t.Errorf("expected custom URL, got %s", cfg.OllamaURL)
	}
	if cfg.DefaultModel != "mistral" {
		t.Errorf("expected mistral model, got %s", cfg.DefaultModel)
	}

	// Test environment variable override
	_ = os.Setenv("OLLAMA_HOST", "http://env-override:11434")
	defer os.Unsetenv("OLLAMA_HOST")

	cfg2, _ := LoadConfig(tempDir)
	if cfg2.OllamaURL != "http://env-override:11434" {
		t.Errorf("env var should override config file: got %s", cfg2.OllamaURL)
	}
}
