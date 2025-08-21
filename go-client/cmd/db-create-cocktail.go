package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate/entities/models"
)

var dbCreateCocktailCmd = &cobra.Command{
	Use:   "create-cocktail",
	Short: "Create Coctail class",
	Run:   cmd_db_create_cocktail,
}

func init() {
	dbCmd.AddCommand(dbCreateCocktailCmd)
}

func cmd_db_create_cocktail(cmd *cobra.Command, args []string) {
	client, err := cmd_db_get_client()
	if err != nil {
		log.Fatal(err)
	}

	classObj := &models.Class{
		Class:       "Cocktail",
		Description: "Alcoholic drink",
		Vectorizer:  "none", // bo embedding dostarczamy sami
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

	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Class Cocktail created")
	}
}
