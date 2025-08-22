package wvclient

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type wvclient struct {
	Client *weaviate.Client
}

func New() *wvclient {
	var err error

	wv := &wvclient{}

	wv.Client, err = get_client("http", "weaviate:8080")
	if err != nil {
		log.Fatal("wvclient error:", err)
	}

	return wv
}

func get_client(scheme string, host string) (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Scheme: scheme,
		Host:   host,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	_, err = client.Misc().ReadyChecker().Do(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}
