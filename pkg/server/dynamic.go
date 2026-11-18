package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/cert"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
)

type certSource interface {
	Certificates() (certs []*cmData.GetCertificateContentResponse, ready bool)
}

// DynamicConfigHandler serves Traefik dynamic configuration.
type DynamicConfigHandler struct {
	certFetcher certSource
}

// NewDynamicConfigHandler builds a new DynamicConfigHandler.
func NewDynamicConfigHandler(f *cert.Fetcher) *DynamicConfigHandler {
	return &DynamicConfigHandler{
		certFetcher: f,
	}
}

type Tls struct {
	Certificates []Certificate `json:"certificates"`
}

type Certificate struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// Get handles an HTTP request to the /config endpoint.
func (d *DynamicConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	certs, ready := d.certFetcher.Certificates()
	if !ready {
		// Traefik keeps its previous configuration on a non-200 response
		http.Error(w, "certificates are not loaded yet", http.StatusServiceUnavailable)
		return
	}

	resp := map[string]any{}
	if len(certs) > 0 {
		tls := Tls{Certificates: make([]Certificate, 0, len(certs))}
		for _, c := range certs {
			tls.Certificates = append(tls.Certificates, Certificate{
				CertFile: strings.Join(c.CertificateChain, "\n"),
				KeyFile:  c.PrivateKey,
			})
		}
		resp["tls"] = tls
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
