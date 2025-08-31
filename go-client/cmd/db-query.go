package cmd

import (
	"context"
	"fmt"
	"go-client/lib/appconfig"
	"go-client/lib/appdebug"
	"go-client/lib/cocktail"
	"go-client/lib/tools"
	"log"

	"github.com/spf13/cobra"
)

var dbQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query Coctail class",
	Run:   cmd_db_query,
}

func init() {
	dbCmd.AddCommand(dbQueryCmd)
}

func cmd_db_query(cmd *cobra.Command, args []string) {
	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	searchText := "Sweet exotic"
	fmt.Println("Search text ", searchText)

	result, err := cr.GetListByNearText(searchText, 3)
	if err != nil {
		log.Fatal(err)
	}

	appdebug.PrettyPrint(result)

}
