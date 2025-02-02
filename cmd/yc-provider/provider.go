package main

import (
	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/cert"
	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"
	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/logger"
	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/server"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewConfiguration,
			logger.NewLogger,
			server.NewDynamicConfigHandler,
			server.NewServeMux,
			server.NewLogger,
			cert.NewFetcher,
			cert.NewSDK,
		),
		fx.Invoke(
			server.NewHTTPServer,
		),
	).Run()
}
