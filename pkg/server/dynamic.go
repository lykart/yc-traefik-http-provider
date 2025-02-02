package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/cert"
)

// DynamicConfigHandler is an http.HandlerFunc that copies its request body
// back to the response.
type DynamicConfigHandler struct {
	certFetcher *cert.Fetcher
}

// NewDynamicConfigHandler builds a new DynamicConfigHandler.
func NewDynamicConfigHandler(f *cert.Fetcher) *DynamicConfigHandler {
	return &DynamicConfigHandler{
		certFetcher: f,
	}
}

type response map[string]interface{}

type Tls struct {
	Certificates []Certificate `json:"certificates"`
}

type Certificate struct {
	CertFile string
	KeyFile  string
}

// Get handles an HTTP request to the /config endpoint.
func (d *DynamicConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{}
	if d.certFetcher.CertificateContent != nil {
		resp["tls"] = Tls{
			Certificates: []Certificate{
				{
					CertFile: strings.Join(d.certFetcher.CertificateContent.CertificateChain, "\n"),
					KeyFile:  d.certFetcher.CertificateContent.PrivateKey,
				},
			},
		}
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}
}
