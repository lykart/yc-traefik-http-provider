package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"yc-traefik-http-provider/pkg/config"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

// NewHTTPServer builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, mux *chi.Mux) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(`:%d`, cfg.Port), Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
