package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

// NewServeMux builds a ServeMux that will route requests
// to the given EchoHandler.
func NewServeMux(d *DynamicConfigHandler, logger *httplog.Logger) *chi.Mux {
	// Service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(
		middleware.Recoverer,
		middleware.StripSlashes,
	)
	r.Get("/config", d.ServeHTTP)

	return r
}
