package cmd

import (
	"go-client/lib/aiclient"
	"go-client/lib/appconfig"
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

	httpPort := appconfig.AppCfg.AiChatCfg[chatName].TmpHttpPort

	ai := aiclient.New(cfgFile, sessionId, chatName)
	ch := wschat.New(httpPort, ai.Ask)
	ch.Serve()
}
