package aiclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	a.defineTools()

}

func (a *aiclient) defineTools() {
	fWeather := openai.FunctionDefinition{
		Name:        "get_current_weather",
		Description: "Get the current weather in a given location. Temperature always in celsius.",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"location": {
					Type:        jsonschema.String,
					Description: "The city and state, e.g. San Francisco, CA",
				},
			},
			Required: []string{"location"},
		},
	}
	a.tools = append(
		a.tools,
		openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fWeather,
		},
	)

	fTime := openai.FunctionDefinition{
		Name:        "get_current_time",
		Description: "Get the current time. Response is in YYYY-MM-DD hh:mm:ss format",
		Parameters: jsonschema.Definition{
			Type:       jsonschema.Object,
			Properties: map[string]jsonschema.Definition{},
		},
	}
	a.tools = append(
		a.tools,
		openai.Tool{
			Type:     openai.ToolTypeFunction,
			Function: &fTime,
		},
	)
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
	log.Printf("Sending request to AI: %s", inputMsg)

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: string(inputMsg)})

	response, err := a.request(
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

	return response.Choices[0].Message.Content, nil
}

func (a *aiclient) request(request openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	var err error
	var response openai.ChatCompletionResponse

	response, err = a.client.CreateChatCompletion(a.ctx, request)
	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	respMsg := response.Choices[0].Message
	log.Printf("Response from AI: %s", respMsg.Content)

	a.messages = append(a.messages, respMsg)

	if len(respMsg.ToolCalls) > 0 {
		response, err = a.handleToolCalls(respMsg.ToolCalls)
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
	}

	return response, err
}

func (a *aiclient) handleToolCalls(toolCalls []openai.ToolCall) (openai.ChatCompletionResponse, error) {
	for _, toolCall := range toolCalls {
		log.Printf("Call: %s(#%v)", toolCall.Function.Name, toolCall.Function.Arguments)
		result, err := a.callFunction(toolCall)
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
		log.Printf("Result (%s): %s", toolCall.Function.Name, result)
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    result,
			Name:       toolCall.Function.Name,
			ToolCallID: toolCall.ID,
		})
	}

	log.Printf("Sending request to AI with results(s) from tool(s)")
	// log.Printf("%#v\n", a.messages)
	return a.client.CreateChatCompletion(
		a.ctx,
		openai.ChatCompletionRequest{
			Model:       a.cfg.Model,
			Temperature: a.cfg.Temperature,
			Messages:    a.messages,
		},
	)
}

func (a *aiclient) callFunction(toolCall openai.ToolCall) (string, error) {
	switch toolCall.Function.Name {
	case "get_current_weather":
		return "Sunny and 36 degrees.", nil
	case "get_current_time":
		t := time.Now()
		return t.Format("2006-01-02 15:04:05"), nil
	default:
		return "", fmt.Errorf("Unknown function name:", toolCall.Function.Name)
	}
}
