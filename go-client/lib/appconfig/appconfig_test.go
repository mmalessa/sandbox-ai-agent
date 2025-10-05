package appconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	yamlContent := `
weaviate:
  scheme: http
  host: localhost:8080
chats:
  coordinator:
    model: "qwen3:1.7b"
    temperature: 0.5
    availableFunctions: ["get_current_time"]
    prompt:
      role: "test role"
      instructions: "test instructions"
    tmpHttpPort: 3000
functions:
  get_drink_recipe:
    url: "http://localhost:3002/api/ask"
    description: "recipe"
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	if err := LoadConfig(configPath); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if AppCfg == nil {
		t.Fatal("expected AppCfg to be populated, got nil")
	}

	if AppCfg.Weaviate == nil || AppCfg.Weaviate.Host != "localhost:8080" {
		t.Errorf("unexpected Weaviate config: %+v", AppCfg.Weaviate)
	}

	if AppCfg.AiChatCfg == nil {
		t.Fatal("expected chats to be populated, got nil")
	}
	chat, ok := AppCfg.AiChatCfg["coordinator"]
	if !ok {
		t.Fatal("expected coordinator chat to exist")
	}
	if chat.Model != "qwen3:1.7b" || chat.TmpHttpPort != 3000 {
		t.Errorf("unexpected chat config: %+v", chat)
	}

	if AppCfg.FunctionCfg == nil || AppCfg.FunctionCfg["get_drink_recipe"].Url == "" {
		t.Errorf("unexpected FunctionCfg: %+v", AppCfg.FunctionCfg)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Invalid YAML
	yamlContent := `
weaviate:
  scheme: http
  host: localhost:8080
chats:
  coordinator:
    model: "qwen3:1.7b
`
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	err := LoadConfig(configPath)
	if err == nil {
		t.Fatal("expected YAML parsing error, got nil")
	}
}
