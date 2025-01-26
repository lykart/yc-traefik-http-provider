package main

import (
	"net/http"

	"yc-traefik-http-provider/pkg/config"
	"yc-traefik-http-provider/pkg/fetcher"
	"yc-traefik-http-provider/pkg/server"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			server.NewHTTPServer,
			server.NewDynamicConfigHandler,
			server.NewServeMux,
			fetcher.NewFetcher,
			fetcher.NewSDK,
		),
		fx.Invoke(
			func(*fetcher.Fetcher) {},
			func(*http.Server) {},
		),
	).Run()
}
