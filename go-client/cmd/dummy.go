package cmd

import (
	"encoding/json"
	"fmt"
	"go-client/lib/appconfig"
	"log"

	"github.com/spf13/cobra"
)

var dummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "Run dummy command",
	Run:   cmd_dummy,
}

func init() {
	rootCmd.AddCommand(dummyCmd)
}

func cmd_dummy(cmd *cobra.Command, args []string) {
	log.Println("Dummy command here")
	b, _ := json.MarshalIndent(appconfig.AppCfg, "", "  ")
	fmt.Println(string(b))
}
