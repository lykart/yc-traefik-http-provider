package server

import (
	"github.com/go-chi/chi/v5"
)

// NewServeMux builds a ServeMux that will route requests
// to the given EchoHandler.
func NewServeMux(d *DynamicConfigHandler) *chi.Mux {
	mux := chi.NewMux()
	mux.Get("/config", d.ServeHTTP)
	return mux
}
