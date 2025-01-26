package fetcher

import (
	"context"
	"log"
	"time"

	"yc-traefik-http-provider/pkg/config"

	cmData "github.com/yandex-cloud/go-genproto/yandex/cloud/certificatemanager/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/fx"
)

type Fetcher struct {
	sdk *ycsdk.SDK
	cfg *config.Config

	CertificateContent *cmData.GetCertificateContentResponse
}

func (f *Fetcher) GetContent(ctx context.Context) {
	c, err := f.sdk.CertificatesData().CertificateContent().Get(ctx, &cmData.GetCertificateContentRequest{
		CertificateId: f.cfg.CertificateId,
	})
	if err != nil {
		log.Printf("failed to get certificate content: %v", err)
		//return err
	}
	log.Printf("got certificate %s", c.CertificateId)
	f.CertificateContent = c
}

func NewFetcher(lc fx.Lifecycle, sdk *ycsdk.SDK, cfg *config.Config) *Fetcher {
	ctx, cancel := context.WithCancel(context.Background())

	f := &Fetcher{
		sdk: sdk,
		cfg: cfg,
	}
	timer := time.NewTicker(30 * time.Second)
	f.GetContent(ctx)
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				for {
					select {
					case <-timer.C:
						f.GetContent(ctx)
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
	return f
}
