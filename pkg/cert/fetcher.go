package cert

import (
	"context"
	"log"
	"time"

	"github.com/nikolaymatrosov/yc-traefik-http-provider/pkg/config"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/fx"
)

type Fetcher struct {
	sdk *ycsdk.SDK
	cfg *config.Configuration

	CertificateContent *cmData.GetCertificateContentResponse
}

func (f *Fetcher) getContent(ctx context.Context) {
	c, err := f.sdk.CertificatesData().CertificateContent().Get(ctx, &cmData.GetCertificateContentRequest{
		CertificateId: f.cfg.Features.CertificateManager.CertificateId,
	})
	if err != nil {
		log.Printf("failed to get certificate content: %v", err)
	}
	f.CertificateContent = c
}

func NewFetcher(lc fx.Lifecycle, sdk *ycsdk.SDK, cfg *config.Configuration) *Fetcher {
	f := &Fetcher{
		sdk: sdk,
		cfg: cfg,
	}

	if f.cfg.Features.CertificateManager.Enabled && f.cfg.Features.CertificateManager.CertificateId != "" {
		ctx, cancel := context.WithCancel(context.Background())
		timer := time.NewTicker(30 * time.Second)
		f.getContent(ctx)
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-timer.C:
							f.getContent(ctx)
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
