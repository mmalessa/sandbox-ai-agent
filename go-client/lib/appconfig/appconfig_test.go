package appconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("cannot create temp config file: %v", err)
	}
	return tmpFile
}

func TestLoadConfig_Success(t *testing.T) {
	t.Cleanup(func() { AppCfg = nil })
	yamlContent := `
chats:
  default:
    model: "gpt-4"
    temperature: 0.7
    prompt:
      role: "assistant"
      context: "general"
      examples: "Q: Hi A: Hello"
      task: "answer questions"
      instructions: "be polite"
    availableFunctions: ["search", "math"]
    tmpHttpPort: 8080
functions:
  calc:
    url: "http://localhost:8081/calc"
    description: "Simple calculator"
weaviate:
  scheme: "http"
  host: "localhost:8082"
`
	path := createTempConfigFile(t, yamlContent)

	if err := LoadConfig(path); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if AppCfg == nil {
		t.Fatal("AppCfg is nil after LoadConfig")
	}

	// Chats
	cfg, ok := AppCfg.AiChatCfg["default"]
	if !ok {
		t.Fatal("default chat config not found")
	}
	if cfg.Model != "gpt-4" {
		t.Errorf("expected model=gpt-4, got=%s", cfg.Model)
	}
	if cfg.TmpHttpPort != 8080 {
		t.Errorf("expected tmpHttpPort=8080, got=%d", cfg.TmpHttpPort)
	}
	if len(cfg.AvailableFunctions) != 2 {
		t.Errorf("expected 2 functions, got=%d", len(cfg.AvailableFunctions))
	}

	// Functions
	fn, ok := AppCfg.FunctionCfg["calc"]
	if !ok {
		t.Fatal("calc function config not found")
	}
	if fn.Url != "http://localhost:8081/calc" {
		t.Errorf("expected url=http://localhost:8081/calc, got=%s", fn.Url)
	}

	// Weaviate
	if AppCfg.Weaviate == nil {
		t.Fatal("weaviate config is nil")
	}
	if AppCfg.Weaviate.Scheme != "http" || AppCfg.Weaviate.Host != "localhost:8082" {
		t.Errorf("unexpected weaviate config: %+v", AppCfg.Weaviate)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	err := LoadConfig("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
	if AppCfg != nil {
		t.Error("AppCfg should remain nil on error")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	invalidYAML := `
chats:
  default:
    model: gpt-4
    temperature: not-a-float
`
	path := createTempConfigFile(t, invalidYAML)

	err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
