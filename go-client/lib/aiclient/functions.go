package aiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"go-client/lib/appconfig"
	"go-client/lib/cocktail"
	"go-client/lib/tools"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type functionDef struct {
	definition openai.FunctionDefinition
	callFn     func(toolCall openai.ToolCall, sessionId string) (string, error)
}

var toolFunctions []functionDef = []functionDef{
	{
		definition: openai.FunctionDefinition{
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
		},
		callFn: func(toolCall openai.ToolCall, sessionId string) (string, error) {
			return "{\"temperature_celsius\": 23.5,\"pressure_hpa\": 1013,\"conditions\": \"Partly cloudy\",\"humidity_percent\": 65}", nil
		},
	},
	{
		definition: openai.FunctionDefinition{
			Name:        "get_current_time",
			Description: "Get the current time. Response is in YYYY-MM-DD hh:mm:ss format",
			Parameters: jsonschema.Definition{
				Type:       jsonschema.Object,
				Properties: map[string]jsonschema.Definition{},
			},
		},
		callFn: func(toolCall openai.ToolCall, sessionId string) (string, error) {
			t := time.Now()
			return t.Format("2006-01-02 15:04:05"), nil
		},
	},
	{
		definition: openai.FunctionDefinition{
			Name:        "cocktail_list",
			Description: "Get list of alcohol cocktails by entering user description",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"user_description": {
						Type:        jsonschema.String,
						Description: "User description",
					},
				},
			},
		},
		callFn: GetCocktailList,
	},
	{
		definition: openai.FunctionDefinition{
			Name:        "cocktail_recipe",
			Description: "Get recipe for a cocktail with a user specified name",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"cocktail_name": {
						Type:        jsonschema.String,
						Description: "Cocktail name",
					},
				},
			},
		},
		callFn: GetCocktailIstructions,
	},
}

func GetCocktailList(toolCall openai.ToolCall, sessionId string) (string, error) {

	// Find user description
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		log.Fatal(err)
	}
	userRequest := args["user_description"].(string)

	// Weaviate Client
	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	// Query
	limit := 5
	cocktails, err := cr.GetListByNearText(userRequest, limit)
	if err != nil {
		return "", err
	}

	// Build string response
	var builder strings.Builder
	for _, c := range cocktails {
		builder.WriteString(fmt.Sprintf("%s (%s)\n", c.Name, c.Ingredients))
	}

	return builder.String(), nil
}

func GetCocktailIstructions(toolCall openai.ToolCall, sessionId string) (string, error) {

	// Find user description
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		log.Fatal(err)
	}
	cocktailName := args["cocktail_name"].(string)

	// Weaviate Client
	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	// Query
	cocktail, err := cr.GetByCocktailName(cocktailName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Name: %s\nIngredients: %s\nPreparation: %s\n", cocktail.Name, cocktail.Ingredients, cocktail.Preparation), nil
}
