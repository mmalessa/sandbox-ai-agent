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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := appconfig.LoadConfig(cfgFile); err != nil {
			log.Fatal(err)
			return err
		}
		log.Printf("Check if configuration for chat \"%s\" exists", chatName)
		if _, ok := appconfig.AppCfg.AiChatCfg[chatName]; !ok {
			err := fmt.Errorf("Configuration for chat \"%s\" not found", chatName)
			log.Fatal(err)
			return err
		}
		return nil
	},
}

var cfgFile string
var chatName string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "Config file path")
	rootCmd.PersistentFlags().StringVarP(&chatName, "chat", "", "coordinator", "Chat name in config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
