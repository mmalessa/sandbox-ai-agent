package tools

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func GetWeaviateClient(scheme string, host string) *weaviate.Client {
	cfg := weaviate.Config{
		Scheme: scheme,
		Host:   host,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return client
}
