package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type chatConfig struct {
	Model       string  `yaml:"model"`
	Temperature float32 `yaml:"temperature"`
	Messages    []struct {
		Role    string `yaml:"role"`
		Content string `yaml:"content"`
	} `yaml:"messages"`
}

func loadChatConfig(path string) (*chatConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", path, err)
	}

	var cfg chatConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("YAML parsing error: %w", err)
	}

	return &cfg, nil
}

func getChatCompletionMessages(chatConfig *chatConfig) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, len(chatConfig.Messages))
	for i, m := range chatConfig.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	return messages
}

func main() {

	baseURL := os.Getenv("OPENAI_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}

	apiToken := os.Getenv("OPENAI_API_TOKEN")
	if apiToken == "" {
		apiToken = "DefaultToken"
	}

	cfg, err := loadChatConfig("chat.yaml")
	if err != nil {
		log.Fatal("error loading YAML:", err)
	}

	config := openai.DefaultConfig(apiToken)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	messages := getChatCompletionMessages(cfg)

	fmt.Println("Sending request")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       cfg.Model,
			Temperature: cfg.Temperature,
			Messages:    messages,
		},
	)
	if err != nil {
		log.Fatalf("Sending request ERROR: %v", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Println("\n---------------------")
}
