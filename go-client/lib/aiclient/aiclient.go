package aiclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"

	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type chatConfig struct {
	Model              string   `yaml:"model"`
	Temperature        float32  `yaml:"temperature"`
	SystemMessage      string   `yaml:"systemMessage"`
	AvailableFunctions []string `yaml:"availableFunctions"`
}

type aiclient struct {
	ctx       context.Context
	baseURL   string
	apiToken  string
	sessionId string
	cfg       *chatConfig
	client    *openai.Client
	messages  []openai.ChatCompletionMessage
	tools     []openai.Tool
}

func New(cfgFile string, sessionId string) *aiclient {

	log.Printf("New session with ID %s started", sessionId)

	apiToken := os.Getenv("OPENAI_API_TOKEN")
	if apiToken == "" {
		apiToken = "DefaultToken"

	}

	a := &aiclient{
		ctx:       context.Background(),
		baseURL:   os.Getenv("OPENAI_URL"),
		apiToken:  apiToken,
		sessionId: sessionId,
		cfg:       &chatConfig{},
	}
	if err := a.loadChatConfig(cfgFile); err != nil {
		log.Fatal("error loading YAML:", err)
	}

	a.initAiClient()
	return a
}

func (a *aiclient) initAiClient() {
	config := openai.DefaultConfig(a.apiToken)
	if a.baseURL != "" {
		config.BaseURL = a.baseURL
	}
	a.client = openai.NewClientWithConfig(config)

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: "system", Content: a.cfg.SystemMessage})

	a.defineTools()

	log.Printf("Available functions: %s", a.cfg.AvailableFunctions)
}

func (a *aiclient) loadChatConfig(path string) error {
	log.Println("Load config file:", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, a.cfg); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	log.Println("Config file loaded")
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
		log.Println("Start using tools")
		response, err = a.handleToolCalls(respMsg.ToolCalls)
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
	}

	return response, err
}

func (a *aiclient) handleToolCalls(toolCalls []openai.ToolCall) (openai.ChatCompletionResponse, error) {
	for _, toolCall := range toolCalls {
		log.Printf("Call function: %s(#%v)", toolCall.Function.Name, toolCall.Function.Arguments)
		result, err := a.callFunction(toolCall)
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
		log.Printf("Function result (%s): %s", toolCall.Function.Name, result)
		a.messages = append(a.messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    result,
			Name:       toolCall.Function.Name,
			ToolCallID: toolCall.ID,
		})
	}

	log.Printf("Sending request to AI with results(s) from tool(s)")

	response, err := a.client.CreateChatCompletion(
		a.ctx,
		openai.ChatCompletionRequest{
			Model:       a.cfg.Model,
			Temperature: a.cfg.Temperature,
			Messages:    a.messages,
		},
	)

	log.Printf("Response from AI with result(s) from tool(s): %s\n", response.Choices[0].Message.Content)

	return response, err
}

func (a *aiclient) defineTools() {
	for _, f := range toolFunctions {
		if slices.Contains(a.cfg.AvailableFunctions, f.definition.Name) {
			fmt.Println("Available function:", f.definition.Name)
			a.tools = append(
				a.tools,
				openai.Tool{
					Type:     openai.ToolTypeFunction,
					Function: &f.definition,
				},
			)
		}
	}
}

func (a *aiclient) callFunction(toolCall openai.ToolCall) (string, error) {
	for _, f := range toolFunctions {
		if f.definition.Name == toolCall.Function.Name {
			return f.callFn(toolCall, a.sessionId)
		}
	}
	return "", fmt.Errorf("unknown function name: %s", toolCall.Function.Name)
}
