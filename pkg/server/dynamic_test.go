package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
)

type fakeSource struct {
	certs []*cmData.GetCertificateContentResponse
	ready bool
}

func (s fakeSource) Certificates() ([]*cmData.GetCertificateContentResponse, bool) {
	return s.certs, s.ready
}

func serve(s fakeSource) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h := &DynamicConfigHandler{certFetcher: s}
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/config", nil))
	return rec
}

func TestConfigNotReady(t *testing.T) {
	rec := serve(fakeSource{ready: false})
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("got %d, want 503", rec.Code)
	}
}

func TestConfigNoCertificates(t *testing.T) {
	rec := serve(fakeSource{ready: true})
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "{}" {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}
}

func TestConfigCertificates(t *testing.T) {
	rec := serve(fakeSource{ready: true, certs: []*cmData.GetCertificateContentResponse{
		{CertificateChain: []string{"leaf1", "ca1"}, PrivateKey: "key1"},
		{CertificateChain: []string{"leaf2"}, PrivateKey: "key2"},
	}})
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("content type %q", ct)
	}
	want := `{"tls":{"certificates":[{"certFile":"leaf1\nca1","keyFile":"key1"},{"certFile":"leaf2","keyFile":"key2"}]}}`
	if got := rec.Body.String(); got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
