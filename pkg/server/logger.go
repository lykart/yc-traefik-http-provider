package server

import (
	"time"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"
	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/logger"

	"github.com/go-chi/httplog/v2"
)

func NewLogger(c *config.Configuration) *httplog.Logger {
	return httplog.NewLogger("http", httplog.Options{
		JSON:             true,
		LogLevel:         logger.LogLevel(c.LogLevel),
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})
}
