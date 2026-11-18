package cert

import (
	"context"
	"errors"
	"testing"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
	"google.golang.org/grpc"
)

type fakeGetter struct {
	keys map[string]string
	errs map[string]error
}

func (g *fakeGetter) Get(_ context.Context, in *cmData.GetCertificateContentRequest, _ ...grpc.CallOption) (*cmData.GetCertificateContentResponse, error) {
	if err := g.errs[in.CertificateId]; err != nil {
		return nil, err
	}
	return &cmData.GetCertificateContentResponse{
		CertificateId: in.CertificateId,
		PrivateKey:    g.keys[in.CertificateId],
	}, nil
}

func keys(t *testing.T, f *Fetcher) []string {
	t.Helper()
	certs, ready := f.Certificates()
	if !ready {
		t.Fatal("fetcher is not ready")
	}
	var out []string
	for _, c := range certs {
		out = append(out, c.PrivateKey)
	}
	return out
}

func TestFetcherNoIDs(t *testing.T) {
	f := newFetcher(&fakeGetter{}, nil)
	if certs, ready := f.Certificates(); !ready || len(certs) != 0 {
		t.Fatalf("got %v, ready=%v", certs, ready)
	}
}

func TestFetcherNotReadyUntilAllLoaded(t *testing.T) {
	g := &fakeGetter{
		keys: map[string]string{"a": "ka", "b": "kb"},
		errs: map[string]error{"b": errors.New("unavailable")},
	}
	f := newFetcher(g, []string{"a", "b"})

	if _, ready := f.Certificates(); ready {
		t.Fatal("ready before first fetch")
	}
	f.fetchAll(context.Background())
	if _, ready := f.Certificates(); ready {
		t.Fatal("ready while one certificate was never fetched")
	}

	g.errs = nil
	f.fetchAll(context.Background())
	got := keys(t, f)
	if len(got) != 2 || got[0] != "ka" || got[1] != "kb" {
		t.Fatalf("got %v", got)
	}
}

func TestFetcherKeepsCacheOnError(t *testing.T) {
	g := &fakeGetter{keys: map[string]string{"a": "ka1", "b": "kb1"}}
	f := newFetcher(g, []string{"a", "b"})
	f.fetchAll(context.Background())

	g.keys = map[string]string{"a": "ka2", "b": "kb2"}
	g.errs = map[string]error{"a": errors.New("unavailable")}
	f.fetchAll(context.Background())

	got := keys(t, f)
	if len(got) != 2 || got[0] != "ka1" || got[1] != "kb2" {
		t.Fatalf("got %v, want [ka1 kb2]", got)
	}
}
