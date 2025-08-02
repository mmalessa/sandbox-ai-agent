package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
	"gopkg.in/yaml.v3"
)

type ChatConfig struct {
	Model    string        `yaml:"model"`
	Messages []api.Message `yaml:"messages"`
}

func loadChatConfig(path string) (*ChatConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ChatConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := loadChatConfig("chat.yaml")
	if err != nil {
		log.Fatal("Błąd wczytywania YAML:", err)
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    cfg.Model,
		Messages: cfg.Messages,
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n---------------------")
}
