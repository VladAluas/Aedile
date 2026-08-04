package platform

import (
	"context"

	"github.com/VladAluas/Aedile/internal/metadata"
)

func (p *Platform) Seed(ctx context.Context) error {

	conn, err := p.DB()
	if err != nil {
		return err
	}

	cfg, err := metadata.Load()
	if err != nil {
		return err
	}

	return metadata.Seed(ctx, conn, cfg)
}
