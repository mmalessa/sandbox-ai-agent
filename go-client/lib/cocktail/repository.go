package cocktail

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type cocktailRepository struct {
	client *weaviate.Client
	ctx    *context.Context
}

func NewRepository(client *weaviate.Client, ctx *context.Context) *cocktailRepository {
	return &cocktailRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *cocktailRepository) InitClass() error {
	err := r.client.Schema().ClassCreator().WithClass(&CocktailClass).Do(*r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *cocktailRepository) ClearClass() error {
	err := r.client.Schema().ClassDeleter().WithClassName(CocktailClassName).Do(*r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *cocktailRepository) Save(c Cocktail) error {
	_, err := r.client.Data().Creator().
		WithClassName(CocktailClassName).
		WithProperties(map[string]interface{}{
			"name":        c.Name,
			"ingredients": c.Ingredients,
			"preparation": c.Preparation,
		}).
		Do(context.Background())

	if err != nil {
		return err
	}
	return nil
}

// Queries
func (r *cocktailRepository) GetListByNearText(text string, limit int) ([]Cocktail, error) {
	fields := []graphql.Field{
		{Name: "name"},
		{Name: "ingredients"},
		{Name: "preparation"},
	}

	nearText := r.client.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{text})

	response, err := r.client.GraphQL().Get().
		WithClassName(CocktailClassName).
		WithFields(fields...).
		WithNearText(nearText).
		WithLimit(limit).
		Do(*r.ctx)
	if err != nil {
		log.Fatal("WV Fatal", err)
	}

	return r.buildResult(response)
}

func (r *cocktailRepository) GetByCocktailName(name string) (Cocktail, error) {
	fields := []graphql.Field{
		{Name: "name"},
		{Name: "ingredients"},
		{Name: "preparation"},
	}

	where := filters.Where().
		WithPath([]string{"name"}).
		WithOperator(filters.Equal).
		WithValueString(name)

	response, err := r.client.GraphQL().Get().
		WithClassName(CocktailClassName).
		WithFields(fields...).
		WithWhere(where).
		Do(*r.ctx)
	if err != nil {
		log.Fatal("WV Fatal", err)
	}

	result, err := r.buildResult(response)
	if err != nil {
		return Cocktail{}, err
	}

	if len(result) > 1 {
		return Cocktail{}, fmt.Errorf("More than one result found")
	}

	if len(result) == 0 {
		return Cocktail{}, nil
	}

	return result[0], nil
}

func (r *cocktailRepository) buildResult(result *models.GraphQLResponse) ([]Cocktail, error) {
	cocktails := []Cocktail{}

	if len(result.Errors) > 0 {
		for _, gqlErr := range result.Errors {
			log.Printf("Błąd GraphQL: %s", gqlErr.Message)
		}
		log.Fatal("END")
	}

	// build response
	getData, ok := result.Data["Get"].(map[string]interface{})
	if !ok {
		return cocktails, fmt.Errorf("no Get field in GQL response")
	}

	getCocktails, ok := getData["Cocktail"].([]interface{})
	if !ok {
		return cocktails, fmt.Errorf("no Cocktail field in GQL response")
	}

	for _, c := range getCocktails {
		cocktail := c.(map[string]interface{})
		cocktails = append(
			cocktails,
			Cocktail{
				Name:        cocktail["name"].(string),
				Ingredients: cocktail["ingredients"].(string),
				Preparation: cocktail["preparation"].(string),
			},
		)
	}
	return cocktails, nil
}
