package fetcher

import (
	"context"

	"yc-traefik-http-provider/pkg/config"

	ycsdk "github.com/yandex-cloud/go-sdk"
)

func NewSDK(config *config.Config) (*ycsdk.SDK, error) {
	var creds ycsdk.Credentials
	if config.SAKey == nil {
		creds = ycsdk.InstanceServiceAccount()
	} else {
		c, err := ycsdk.ServiceAccountKey(config.SAKey)
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
