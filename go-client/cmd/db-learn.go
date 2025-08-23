package cmd

import (
	"context"
	"go-client/lib/appconfig"
	"go-client/lib/cocktail"
	"go-client/lib/tools"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
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
	var cocktails []cocktail.Cocktail

	path := "data/cocktails.csv"

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := gocsv.UnmarshalFile(f, &cocktails); err != nil {
		log.Fatal(err)
	}

	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	for _, c := range cocktails {
		if err := cr.Save(c); err != nil {
			log.Fatal(err)
		}
		log.Printf("Added: %s", c.Name)
	}
}
