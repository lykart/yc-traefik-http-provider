package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"yc-traefik-http-provider/pkg/fetcher"
)

// DynamicConfigHandler is an http.HandlerFunc that copies its request body
// back to the response.
type DynamicConfigHandler struct {
	fetcher *fetcher.Fetcher
}

// NewDynamicConfigHandler builds a new DynamicConfigHandler.
func NewDynamicConfigHandler(f *fetcher.Fetcher) *DynamicConfigHandler {
	return &DynamicConfigHandler{
		fetcher: f,
	}
}

type response struct {
	Tls Tls `json:"tls"`
}

type Tls struct {
	Certificates []Certificate `json:"certificates"`
}

type Certificate struct {
	CertFile string
	KeyFile  string
}

// Get handles an HTTP request to the /config endpoint.
func (d *DynamicConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d.fetcher.CertificateContent == nil {
		d.fetcher.GetContent(r.Context())
	}

	resp := response{
		Tls{
			Certificates: []Certificate{
				{
					CertFile: strings.Join(d.fetcher.CertificateContent.CertificateChain, "\n"),
					KeyFile:  d.fetcher.CertificateContent.PrivateKey,
				},
			},
		},
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

	//if _, err := io.Copy(w, r.Body); err != nil {
	//	fmt.Fprintln(os.Stderr, "Failed to handle request:", err)
	//}
}
