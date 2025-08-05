package aiclient

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
	System      string  `yaml:"system"`
}

type aiclient struct {
	ctx      context.Context
	baseURL  string
	apiToken string
	cfg      *chatConfig
	client   *openai.Client
	messages []openai.ChatCompletionMessage
}

func New() *aiclient {
	apiToken := os.Getenv("OPENAI_API_TOKEN")
	if apiToken == "" {
		apiToken = "DefaultToken"

	}

	a := &aiclient{
		ctx:      context.Background(),
		baseURL:  os.Getenv("OPENAI_URL"),
		apiToken: apiToken,
	}
	a.init()
	return a
}

// FIXME
func (a *aiclient) init() {

	if err := a.loadChatConfig("chat.yaml"); err != nil {
		log.Fatal("error loading YAML:", err)
	}

	// config := openai.DefaultConfig(a.apiToken)
	// if a.baseURL != "" {
	// 	config.BaseURL = a.baseURL
	// }
	// a.client = openai.NewClientWithConfig(config)

	// messages = append(messages, openai.ChatCompletionMessage{Role: "system", Content: cfg.System})

	// some tools
	// describe the function & its inputs
	// params := jsonschema.Definition{
	// 	Type: jsonschema.Object,
	// 	Properties: map[string]jsonschema.Definition{
	// 		"location": {
	// 			Type:        jsonschema.String,
	// 			Description: "The city and state, e.g. San Francisco, CA",
	// 		},
	// 		"unit": {
	// 			Type: jsonschema.String,
	// 			Enum: []string{"celsius", "fahrenheit"},
	// 		},
	// 	},
	// 	Required: []string{"location"},
	// }
	// f := openai.FunctionDefinition{
	// 	Name:        "get_current_weather",
	// 	Description: "Get the current weather in a given location. Only use it when the question is about the weather.",
	// 	Parameters:  params,
	// }
	// t := openai.Tool{
	// 	Type:     openai.ToolTypeFunction,
	// 	Function: &f,
	// }

	// tools = append(tools, t)
}

func (a *aiclient) loadChatConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, a.cfg); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	return nil
}

func (a *aiclient) Send(inputMsg string) (string, error) {
	// response logic
	log.Printf("Sending to AI: %s", inputMsg)

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: "user", Content: string(inputMsg)})

	resp, err := a.client.CreateChatCompletion(
		a.ctx,
		openai.ChatCompletionRequest{
			Model:       a.cfg.Model,
			Temperature: a.cfg.Temperature,
			Messages:    a.messages,
			// Tools:       tools,
			// ToolChoice:  "auto",
		},
	)
	if err != nil {
		log.Fatalf("Sending request ERROR: %v", err)
		return "", err
	}

	respMsg := resp.Choices[0].Message
	a.messages = append(a.messages, respMsg)

	return respMsg.Content, nil
}
