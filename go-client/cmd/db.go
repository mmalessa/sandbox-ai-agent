package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "DB tools",
}

func init() {
	rootCmd.AddCommand(dbCmd)
}

func cmd_db_get_client() (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Scheme: "http",
		Host:   "weaviate:8080",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	ready, err := client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		return nil, err
	}
	log.Printf("%#v", ready)
	return client, nil
}
