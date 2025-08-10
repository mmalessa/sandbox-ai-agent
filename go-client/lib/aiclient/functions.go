package aiclient

import (
	"encoding/json"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
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
			return "Słonecznie, temperatura 29 stopni Celsjusza", nil
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
}
