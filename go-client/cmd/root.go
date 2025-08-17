package cmd

import (
	"fmt"
	"log"
	"os"

	"go-client/lib/appconfig"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "My AI sandbox",
	Long:  "A simple of AI sandbox",
}

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "agent-main.yaml", "Config file path")

	if err := appconfig.LoadConfig(cfgFile); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
