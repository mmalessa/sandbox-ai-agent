package cmd

import (
	"context"
	"fmt"
	"go-client/lib/wvclient"
	"log"

	"github.com/spf13/cobra"
)

var dbClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete Coctail class",
	Run:   cmd_db_clear,
}

func init() {
	dbCmd.AddCommand(dbClearCmd)
}

func cmd_db_clear(cmd *cobra.Command, args []string) {
	wv := wvclient.New()

	err := wv.Client.Schema().ClassDeleter().WithClassName("Cocktail").Do(context.Background())
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Class Cocktail deleted")
	}
}
