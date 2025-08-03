package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"

	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

var client *openai.Client
var cfg *chatConfig
var messages []openai.ChatCompletionMessage

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

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", msg)
		messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: string(msg)})

		// response logic
		log.Println("Sending to AI")
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
		response := resp.Choices[0].Message.Content
		log.Println(response)

		messages = append(messages, openai.ChatCompletionMessage{Role: "assistant", Content: response})

		// send response to web client
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
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
	// ai client
	baseURL := os.Getenv("OPENAI_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}

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
	config.BaseURL = baseURL
	client = openai.NewClientWithConfig(config)

	messages = append(messages, openai.ChatCompletionMessage{Role: "system", Content: cfg.System})
}

func main() {

	initAiClient()

	// http/websocket server
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsHandler)

	log.Println("Serwer startuje na :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
