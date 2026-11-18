package cert

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

const fetchTimeout = 30 * time.Second

type contentGetter interface {
	Get(ctx context.Context, in *cmData.GetCertificateContentRequest, opts ...grpc.CallOption) (*cmData.GetCertificateContentResponse, error)
}

type Fetcher struct {
	getter contentGetter
	ids    []string

	mu    sync.RWMutex
	cache map[string]*cmData.GetCertificateContentResponse
}

func newFetcher(getter contentGetter, ids []string) *Fetcher {
	return &Fetcher{
		getter: getter,
		ids:    ids,
		cache:  make(map[string]*cmData.GetCertificateContentResponse, len(ids)),
	}
}

// Certificates returns cached contents in configured order.
// ready is false until every configured certificate has been fetched at least once.
func (f *Fetcher) Certificates() (certs []*cmData.GetCertificateContentResponse, ready bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, id := range f.ids {
		c, ok := f.cache[id]
		if !ok {
			return nil, false
		}
		certs = append(certs, c)
	}
	return certs, true
}

func (f *Fetcher) fetchAll(ctx context.Context) {
	for _, id := range f.ids {
		callCtx, cancel := context.WithTimeout(ctx, fetchTimeout)
		c, err := f.getter.Get(callCtx, &cmData.GetCertificateContentRequest{
			CertificateId: id,
		})
		cancel()
		if err != nil {
			// keep the previously cached content
			log.Printf("failed to get certificate %s content: %v", id, err)
			continue
		}
		f.mu.Lock()
		f.cache[id] = c
		f.mu.Unlock()
	}
}

func NewFetcher(lc fx.Lifecycle, sdk *ycsdk.SDK, cfg *config.Configuration) *Fetcher {
	var ids []string
	if cfg.Features.CertificateManager.Enabled {
		ids = cfg.Features.CertificateManager.IDs()
	}
	f := newFetcher(sdk.CertificatesData().CertificateContent(), ids)

	if len(ids) > 0 {
		ctx, cancel := context.WithCancel(context.Background())
		timer := time.NewTicker(time.Duration(cfg.CertFetchIntervalSeconds) * time.Second)
		f.fetchAll(ctx)
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-timer.C:
							f.fetchAll(ctx)
						case <-ctx.Done():
							return
						}
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Printf("stopping fetcher")
				timer.Stop()
				cancel()
				return nil
			},
		})
	}
	return f
}
