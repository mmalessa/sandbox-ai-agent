package cmd

import (
	"context"
	"go-client/lib/appconfig"
	"go-client/lib/cocktail"
	"go-client/lib/tools"
	"log"

	"github.com/spf13/cobra"
)

var dbInitcmd = &cobra.Command{
	Use:   "init",
	Short: "Create Cocktail class",
	Run:   cmd_db_init,
}

func init() {
	dbCmd.AddCommand(dbInitcmd)
}

func cmd_db_init(cmd *cobra.Command, args []string) {
	wvc := tools.GetWeaviateClient(appconfig.AppCfg.Weaviate.Scheme, appconfig.AppCfg.Weaviate.Host)
	ctx := context.Background()
	cr := cocktail.NewRepository(wvc, &ctx)

	err := cr.InitClass()
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Class %s created", cocktail.CocktailClassName)
	}
}
