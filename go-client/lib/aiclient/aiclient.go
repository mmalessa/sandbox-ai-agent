package aiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-client/lib/appconfig"
	"go-client/lib/httptools"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type aiclient struct {
	ctx       context.Context
	baseURL   string
	apiToken  string
	sessionId string
	client    *openai.Client
	messages  []openai.ChatCompletionMessage
	tools     []openai.Tool
	cfg       *appconfig.AiChatConfig
}

func New(cfgFile string, sessionId string, chatName string) *aiclient {

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
		cfg:       appconfig.AppCfg.AiChatCfg[chatName],
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

	a.messages = append(a.messages, openai.ChatCompletionMessage{Role: "system", Content: a.cfg.PromptRole})

	a.defineTools()

	// log.Printf("Available functions: %s", a.cfg.AvailableFunctions)
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

	// built-in functions
	for _, f := range toolFunctions {
		if slices.Contains(a.cfg.AvailableFunctions, f.definition.Name) {
			log.Println("Available function:", f.definition.Name)
			a.tools = append(
				a.tools,
				openai.Tool{
					Type:     openai.ToolTypeFunction,
					Function: &f.definition,
				},
			)
		}
	}

	// API based functions
	for k, f := range appconfig.AppCfg.FunctionCfg {
		if slices.Contains(a.cfg.AvailableFunctions, k) {
			log.Printf("API based function: %s (%s)", k, f.Url)

			functionDefinition := &openai.FunctionDefinition{
				Name:        k,
				Description: f.Description,
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"request": {
							Type:        jsonschema.String,
							Description: "Request from user",
						},
					},
					Required: []string{"request"},
				},
			}

			a.tools = append(
				a.tools,
				openai.Tool{
					Type:     openai.ToolTypeFunction,
					Function: functionDefinition,
				},
			)
		}
	}

}

func (a *aiclient) callFunction(toolCall openai.ToolCall) (string, error) {
	// build-in functions
	for _, f := range toolFunctions {
		if f.definition.Name == toolCall.Function.Name {
			return f.callFn(toolCall, a.sessionId)
		}
	}

	// API based functions
	for k, f := range appconfig.AppCfg.FunctionCfg {
		if k == toolCall.Function.Name {
			return a.callApiBasedFunction(toolCall, a.sessionId, f)
		}
	}
	return "", fmt.Errorf("unknown function name: %s", toolCall.Function.Name)
}

func (a *aiclient) callApiBasedFunction(toolCall openai.ToolCall, sessionId string, f *appconfig.FunctionConfig) (string, error) {

	// toolCall args
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		log.Fatal(err)
	}
	userRequest := args["request"].(string)

	// request context
	// TODO - get context from database
	requestContext := "Offer only alcoholic drinks\nThe user is a gourmet\n"

	// build request based on template
	tmpl, err := template.New("msg").Parse(f.RequestTemplate)
	if err != nil {
		log.Fatal(err)
	}
	values := map[string]interface{}{
		"Context": requestContext,
		"Request": userRequest,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, values); err != nil {
		log.Fatal(err)
	}

	// HTTP Request
	log.Printf("Function request to: %s", f.Url)
	requestData := httptools.RequestData{Content: buf.String()}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	client := &http.Client{
		Timeout: 3 * time.Minute,
	}

	req, err := http.NewRequest("POST", f.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-correlationId", sessionId)

	resp, err := client.Do(req)
	if err != nil {
		log.Print("client.Do ERROR")
		// dump, _ := httputil.DumpResponse(resp, true)
		// fmt.Printf("RAW RESPONSE:\n%s\n", dump)
		log.Fatal(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print("io.ReadAll ERROR")
		log.Fatal(err)
		return "", err
	}

	log.Printf("Response from %s for function. Status: %s", f.Url, resp.Status)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(body))

	responseBody := string(body)
	log.Printf("Response content: %s", responseBody)

	return responseBody, nil
}

func (a *aiclient) GetEmbeddingOllama(model string, text string) ([]float32, error) {

	resp, err := a.client.CreateEmbeddings(
		context.Background(),
		openai.EmbeddingRequest{
			Model: openai.EmbeddingModel(model),
			Input: text,
		},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding for text: %s", text)
	}

	return resp.Data[0].Embedding, nil
}
