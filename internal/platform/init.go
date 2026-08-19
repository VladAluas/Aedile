// Package platform orchestrates workflows
package platform

import (
	"context"

	"github.com/VladAluas/Aedile/internal/generator"
)

func (p *Platform) Init(ctx context.Context) error {
	if err := generator.GenerateCompose(p.cfg); err != nil {
		return err
	}

	return p.docker.Compose(ctx, "up", "-d")
}
