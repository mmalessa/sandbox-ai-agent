package cmd

import (
	"go-client/lib/aiclient"
	"go-client/lib/wschat"

	"github.com/google/uuid"
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
	sessionId := uuid.NewString()

	ai := aiclient.New(cfgFile, sessionId)
	ch := wschat.New(8000, ai.Ask)
	ch.Serve()
}
