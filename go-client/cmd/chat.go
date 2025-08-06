package cmd

import (
	"go-client/lib/aiclient"
	"go-client/lib/wschat"

	"github.com/spf13/cobra"
)

var chat = &cobra.Command{
	Use:   "chat",
	Short: "Run AI chat",
	Run:   cmd_chat,
}

func init() {
	rootCmd.AddCommand(chat)
}

func cmd_chat(cmd *cobra.Command, args []string) {
	ai := aiclient.New()
	ch := wschat.New(8000, ai.Ask)
	ch.Serve()
}
