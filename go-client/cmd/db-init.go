package cmd

import (
	"context"
	"fmt"
	"go-client/lib/wvclient"
	"log"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate/entities/models"
)

var dbInitcmd = &cobra.Command{
	Use:   "init",
	Short: "Create Coctail class",
	Run:   cmd_db_init,
}

func init() {
	dbCmd.AddCommand(dbInitcmd)
}

func cmd_db_init(cmd *cobra.Command, args []string) {
	wv := wvclient.New()

	classObj := &models.Class{
		Class:       "Cocktail",
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

	err := wv.Client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Class Cocktail created")
	}
}
