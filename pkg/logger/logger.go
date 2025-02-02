package logger

import (
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"
)

func NewLogger(c *config.Configuration) *log.Logger {
	return slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), LogLevel(c.LogLevel))
}

func LogLevel(val string) slog.Level {
	switch strings.ToLower(val) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
