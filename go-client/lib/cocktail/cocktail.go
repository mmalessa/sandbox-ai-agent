package cocktail

import (
	"github.com/weaviate/weaviate/entities/models"
)

var CocktailClassName = "Cocktail"

var CocktailClass = models.Class{
	Class:       CocktailClassName,
	Description: "Alcoholic drink",
	// Vectorizer:  "none", // bo embedding dostarczamy sami
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
		},
	},
	ModuleConfig: map[string]interface{}{
		"text2vec-transformers": map[string]interface{}{},
	},
}

type Cocktail struct {
	Name        string `csv:"Cocktail Name"`
	Ingredients string `csv:"Ingredients"`
	Preparation string `csv:"Preparation"`
}
