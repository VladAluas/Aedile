package openbao

import (
	"context"

	bao "github.com/openbao/openbao/api/v2"
)

// Client Unseal
func (c *Client) Unseal(ctx context.Context, key string) (*bao.SealStatusResponse, error) {
	return c.client.Sys().UnsealWithContext(ctx, key)
}
