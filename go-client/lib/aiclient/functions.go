package aiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"go-client/lib/wvclient"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type functionDef struct {
	definition openai.FunctionDefinition
	callFn     func(toolCall openai.ToolCall, sessionId string) (string, error)
}

// TODO
// - add API request inside callFn
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
			Name:        "recipe_list",
			Description: "Get a list of available cooking recipes. Response is JSON. Format: { {name:\"\", description:\"\"} }",
			Parameters: jsonschema.Definition{
				Type:       jsonschema.Object,
				Properties: map[string]jsonschema.Definition{},
			},
		},
		callFn: func(toolCall openai.ToolCall, sessionId string) (string, error) {

			list := []struct {
				Name        string
				Description string
			}{
				{Name: "Jajecznica", Description: "Tradycyjna jajecznica na maśle"},
				{Name: "Jajecznica na boczku", Description: "Jajecznica zrobiona na boczku"},
				{Name: "Kanapka z szynką", Description: "Prosta kanapka - kromka chleba, masło, szynka"},
				{Name: "Schabowy", Description: "Schabowy w panierce z bułki tartej"},
			}
			j, err := json.Marshal(list)

			return string(j), err
		},
	},
	{
		definition: openai.FunctionDefinition{
			Name:        "recipe_book",
			Description: "Download a cooking recipe by entering its name.",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"name": {
						Type:        jsonschema.String,
						Description: "Recipe name",
					},
				},
			},
		},
		callFn: func(toolCall openai.ToolCall, sessionId string) (string, error) {
			// recipeName :=
			return "Recipe content - TODO", nil
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
}

func GetCocktailList(toolCall openai.ToolCall, sessionId string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		log.Fatal(err)
	}
	userRequest := args["user_description"].(string)

	wv := wvclient.New()
	ctx := context.Background()

	fields := []graphql.Field{
		{Name: "name"},
		{Name: "ingredients"},
	}

	nearText := wv.Client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{userRequest})

	result, err := wv.Client.GraphQL().Get().
		WithClassName("Cocktail").
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(5).
		Do(ctx)
	if err != nil {
		log.Fatal("WV Fatal", err)
		return "", err
	}
	if len(result.Errors) > 0 {
		for _, gqlErr := range result.Errors {
			log.Printf("Błąd GraphQL: %s", gqlErr.Message)
		}
		log.Fatal("END")
	}

	// build response
	getData, ok := result.Data["Get"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no Get field in GQL response")
	}

	cocktails, ok := getData["Cocktail"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no Cocktail field in GQL response")
	}

	var builder strings.Builder
	for _, c := range cocktails {
		cocktail := c.(map[string]interface{})
		name := cocktail["name"].(string)
		ingredients := cocktail["ingredients"].(string)

		builder.WriteString(fmt.Sprintf("%s (%s)\n", name, ingredients))
	}

	r := builder.String()

	fmt.Println(r)

	return r, nil
}
