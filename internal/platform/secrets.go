package platform

import (
	"context"
	"fmt"
	"os"

	"github.com/VladAluas/Aedile/internal/openbao"
)

func (p *Platform) SeedSecrets(ctx context.Context) error {
	token := os.Getenv("OPENBAO_TOKEN")
	if token == "" {
		return fmt.Errorf("OPENBAO_TOKEN is not set")
	}

	client, err := openbao.New(
		p.cfg.OpenBao,
		token,
	)
	if err != nil {
		return err
	}

	return client.Seed(ctx, p.cfg.Secrets)
}








