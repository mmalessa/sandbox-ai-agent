package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var dbDeleteCocktailCmd = &cobra.Command{
	Use:   "delete-cocktail",
	Short: "Delete Coctail class",
	Run:   cmd_db_delete_cocktail,
}

func init() {
	dbCmd.AddCommand(dbDeleteCocktailCmd)
}

func cmd_db_delete_cocktail(cmd *cobra.Command, args []string) {
	client, err := cmd_db_get_client()
	if err != nil {
		log.Fatal(err)
	}

	err = client.Schema().ClassDeleter().WithClassName("Cocktail").Do(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Class Cocktail deleted")
	}
}
