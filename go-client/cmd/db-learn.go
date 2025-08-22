package cmd

import (
	"context"
	"fmt"
	"go-client/lib/wvclient"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var dbLearnCmd = &cobra.Command{
	Use:   "learn",
	Short: "Learn cocktails from CSV to DB",
	Run:   cmd_db_learn,
}

type CsvCocktail struct {
	Name        string `csv:"Cocktail Name"`
	Bartender   string `csv:"Bartender"`
	Company     string `csv:"Bar/Company"`
	Location    string `csv:"Location"`
	Ingredients string `csv:"Ingredients"`
	Garnish     string `csv:"Garnish"`
	Glassware   string `csv:"Glassware"`
	Preparation string `csv:"Preparation"`
	Notes       string `csv:"Notes"`
}

func init() {
	dbCmd.AddCommand(dbLearnCmd)
}

func cmd_db_learn(cmd *cobra.Command, args []string) {
	cocktails, err := cmd_learn_get_cocktails("data/cocktails.csv")
	if err != nil {
		log.Fatal(err)
	}

	wv := wvclient.New()

	// sessionId := uuid.NewString()
	// ai := aiclient.New(cfgFile, sessionId, chatName)

	for _, c := range cocktails {
		// text := cmd_learn_make_cocktail_text(c)
		// embedding, err := ai.GetEmbeddingOllama("nomic-embed-text:latest", text)
		if err != nil {
			log.Fatal(err)
		}
		cmd_learn_store_cocktail_in_weaviate(wv.Client, c)

		log.Printf("Added: %s", c.Name)
	}
}

func cmd_learn_get_cocktails(path string) ([]CsvCocktail, error) {
	var cocktails []CsvCocktail

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, &cocktails); err != nil {
		return nil, err
	}
	return cocktails, nil
}

func cmd_learn_make_cocktail_text(c CsvCocktail) string {
	return fmt.Sprintf(
		"Cocktail: %s\nIngredients: %s\nPreparation: %s",
		c.Name, c.Ingredients, c.Preparation,
	)
}

// func cmd_learn_store_cocktail_in_weaviate(client *weaviate.Client, c CsvCocktail, vector []float32) {
func cmd_learn_store_cocktail_in_weaviate(client *weaviate.Client, c CsvCocktail) {

	_, err := client.Data().Creator().
		WithClassName("Cocktail").
		WithProperties(map[string]interface{}{
			"name":        c.Name,
			"ingredients": c.Ingredients,
			"preparation": c.Preparation,
		}).
		// WithVector(vector).
		Do(context.Background())

	if err != nil {
		fmt.Println("Error:", err)
	}
}
