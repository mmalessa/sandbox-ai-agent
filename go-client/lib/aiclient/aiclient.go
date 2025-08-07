package aiclient

import (
	"context"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
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
	tools    []openai.Tool
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
		cfg:      &chatConfig{},
	}
	a.init()
	return a
}

// FIXME
func (a *aiclient) init() {

	if err := a.loadChatConfig("chat.yaml"); err != nil {
		log.Fatal("error loading YAML:", err)
	}

	config := openai.DefaultConfig(a.apiToken)
	if a.baseURL != "" {
		config.BaseURL = a.baseURL
	}
	a.client = openai.NewClientWithConfig(config)

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: "system", Content: a.cfg.System})

	// some tools
	// describe the function & its inputs
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"location": {
				Type:        jsonschema.String,
				Description: "The city and state, e.g. San Francisco, CA",
			},
			"unit": {
				Type: jsonschema.String,
				Enum: []string{"celsius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
	}
	f := openai.FunctionDefinition{
		Name:        "get_current_weather",
		Description: "Get the current weather in a given location. Only use it when the question is about the weather.",
		Parameters:  params,
	}
	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &f,
	}

	a.tools = append(a.tools, t)
}

func (a *aiclient) loadChatConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, a.cfg); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	_ = data
	return nil
}

func (a *aiclient) Ask(inputMsg string) (string, error) {
	// response logic
	log.Printf("Sending to AI: %s", inputMsg)

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: "user", Content: string(inputMsg)})

	resp, err := a.request(
		openai.ChatCompletionRequest{
			Model:       a.cfg.Model,
			Temperature: a.cfg.Temperature,
			Messages:    a.messages,
			Tools:       a.tools,
			ToolChoice:  "auto",
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

func (a *aiclient) request(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	response, err := a.client.CreateChatCompletion(a.ctx, request)

	// check ToolCalls
	respMsg := response.Choices[0].Message
	if len(respMsg.ToolCalls) > 0 {
		log.Println("Run ToolCalls")
		tollCall := respMsg.ToolCalls[0]
		// for _, toolCall in respMsg.ToolCalls {
		response, err = a.callFunction(tollCall)
		// }
	}

	return response, err
}

func (a *aiclient) callFunction(toolCall openai.ToolCall) (openai.ChatCompletionResponse, error) {
	if toolCall.Function.Name == "get_current_weather" {
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    "Sunny and 36 degrees.",
			Name:       toolCall.Function.Name,
			ToolCallID: toolCall.ID,
		})

		log.Printf(
			"Sending OpenAI our '%v()' function's response and requesting the reply to the original question...\n",
			toolCall.Function.Name,
		)

		return a.client.CreateChatCompletion(
			a.ctx,
			openai.ChatCompletionRequest{
				Model:       a.cfg.Model,
				Temperature: a.cfg.Temperature,
				Messages:    a.messages,
				// Tools:       a.tools,
			},
		)
	}
	return openai.ChatCompletionResponse{}, fmt.Errorf("FIXME error")
}
