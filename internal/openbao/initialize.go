package openbao

import (
	"context"
	"fmt"

	bao "github.com/openbao/openbao/api/v2"
)

// Client Initialize
func (c *Client) Initialize(
	ctx context.Context,
	shares int,
	threshold int,
	pgpKeys []string,
	rootTokenPGPKey string,
) (*bao.InitResponse, error) {
	initialized, err := c.client.Sys().InitStatusWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("check OpenBao initialization: %w", err)
	}

	if initialized {
		return nil, fmt.Errorf("OpenBao is already initialized")
	}

	if shares < 1 || threshold < 1 || threshold > shares {
		return nil, fmt.Errorf("invalid key shares/threshold")
	}

	if len(pgpKeys) > 0 && len(pgpKeys) != shares {
		return nil, fmt.Errorf(
			"provided %d PGP Keys, but %d key shares were requested",
			len(pgpKeys), shares,
		)
	}

	result, err := c.client.Sys().InitWithContext(ctx, &bao.InitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
		PGPKeys:         pgpKeys,
		RootTokenPGPKey: rootTokenPGPKey,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize OpenBao: %w", err)
	}

	return result, nil
}
