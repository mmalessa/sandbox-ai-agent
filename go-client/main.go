package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"gopkg.in/yaml.v3"
)

var client *openai.Client
var cfg *chatConfig
var messages []openai.ChatCompletionMessage
var ctx context.Context
var tools []openai.Tool

type chatConfig struct {
	Model       string  `yaml:"model"`
	Temperature float32 `yaml:"temperature"`
	System      string  `yaml:"system"`
}

// upgrader changes HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // we allow to connect from any source
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade connection HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection")

	for _, m := range messages {
		if m.Role == openai.ChatMessageRoleSystem {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(m.Content)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	for {
		_, inputMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", inputMsg)
		messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: string(inputMsg)})

		// response logic
		log.Printf("Sending to AI: %s", inputMsg)
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       cfg.Model,
				Temperature: cfg.Temperature,
				Messages:    messages,
				Tools:       tools,
				ToolChoice:  "auto",
			},
		)
		if err != nil {
			log.Fatalf("Sending request ERROR: %v", err)
		}

		respMsg := resp.Choices[0].Message

		if len(respMsg.ToolCalls) > 0 {

			log.Printf("ToolCalls length: %d\n", len(respMsg.ToolCalls))
			toolCall := respMsg.ToolCalls[0]

			if toolCall.Function.Name == "get_current_weather" {
				// simulate calling the function & responding to OpenAI

				messages = append(messages, respMsg)

				log.Printf(
					"OpenAI called us back wanting to invoke our function '%v' with params '%v'\n",
					toolCall.Function.Name,
					toolCall.Function.Arguments)

				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    "Sunny and 36 degrees.",
					Name:       toolCall.Function.Name,
					ToolCallID: toolCall.ID,
				})

				log.Printf(
					"Sending OpenAI our '%v()' function's response and requesting the reply to the original question...\n",
					toolCall.Function.Name)

				resp, err = client.CreateChatCompletion(ctx,
					openai.ChatCompletionRequest{
						Model:       cfg.Model,
						Temperature: cfg.Temperature,
						Messages:    messages,
						Tools:       tools,
					},
				)
				if err != nil {
					log.Fatalf("Sending tool request ERROR: %v", err)
				}
				respMsg = resp.Choices[0].Message
			}
		}

		messages = append(messages, respMsg)

		log.Println(respMsg.Content)

		// send response to web client
		if err := conn.WriteMessage(websocket.TextMessage, []byte(respMsg.Content)); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
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

func initAiClient() {
	ctx = context.Background()

	// ai client
	baseURL := os.Getenv("OPENAI_URL")

	apiToken := os.Getenv("OPENAI_API_TOKEN")
	if apiToken == "" {
		apiToken = "DefaultToken"
	}

	var err error
	cfg, err = loadChatConfig("chat.yaml")
	if err != nil {
		log.Fatal("error loading YAML:", err)
	}

	config := openai.DefaultConfig(apiToken)
	if baseURL != "" {
		// baseURL = "http://localhost:11434/v1"
		config.BaseURL = baseURL
	}
	client = openai.NewClientWithConfig(config)

	messages = append(messages, openai.ChatCompletionMessage{Role: "system", Content: cfg.System})

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

	tools = append(tools, t)
}

func main() {

	initAiClient()

	// http/websocket server
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsHandler)

	log.Println("Serwer startuje na :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
