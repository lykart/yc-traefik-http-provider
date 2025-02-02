package cert

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"

	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

func NewSDK(config *config.Configuration, logger *log.Logger) (*ycsdk.SDK, error) {
	var creds ycsdk.Credentials
	if config.SAKeyPath == "" {
		creds = ycsdk.InstanceServiceAccount()
	} else {
		key, err := serviceAccountKeyFromFile(config.SAKeyPath, logger)
		if err != nil {
			return nil, err
		}
		c, err := ycsdk.ServiceAccountKey(key)
		if err != nil {
			return nil, err
		}
		creds = c
	}
	cfg := ycsdk.Config{
		Credentials: creds,
	}
	ctx := context.Background()
	sdk, err := ycsdk.Build(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func serviceAccountKeyFromFile(path string, logger *log.Logger) (*iamkey.Key, error) {
	logger.Printf("SAKey loaded: %s", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var key iamkey.Key
	if err = json.Unmarshal(data, &key); err != nil {
		return nil, err
	}
	return &key, nil
}
