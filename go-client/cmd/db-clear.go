package cmd

import (
	"context"
	"fmt"
	"go-client/lib/appconfig"
	"go-client/lib/cocktail"
	"go-client/lib/tools"
	"log"

	"github.com/spf13/cobra"
)

var dbClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete Cocktail class",
	Run:   cmd_db_clear,
}

func init() {
	dbCmd.AddCommand(dbClearCmd)
}

func cmd_db_clear(cmd *cobra.Command, args []string) {
	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)

	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	err := cr.ClearClass()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Class Cocktail deleted")
	}
}
