package cocktail

import (
	"github.com/weaviate/weaviate/entities/models"
)

var CocktailClassName = "Cocktail"
var VectorizerName = "text2vec-transformers" // "none" - jeśli embedding dostarczamy sami

var CocktailClass = models.Class{
	Class:       CocktailClassName,
	Description: "Alcoholic drink",
	Vectorizer:  VectorizerName,
	Properties: []*models.Property{
		{
			Name:     "name",
			DataType: []string{"text"},
		},
		{
			Name:     "ingredients",
			DataType: []string{"text"},
		},
		{
			Name:     "preparation",
			DataType: []string{"text"},
			ModuleConfig: map[string]interface{}{
				VectorizerName: map[string]interface{}{
					"skip": true, // exclude from embedding
				},
			},
		},
	},
	ModuleConfig: map[string]interface{}{
		VectorizerName: map[string]interface{}{
			"vectorizeClassName":    false,
			"vectorizePropertyName": true,
		},
	},
}

type Cocktail struct {
	Name        string `csv:"Cocktail Name"`
	Ingredients string `csv:"Ingredients"`
	Preparation string `csv:"Preparation"`
}
