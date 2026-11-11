package openbao

import (
	"context"

	bao "github.com/openbao/openbao/api/v2"
)

// Client Status
func (c *Client) Status(ctx context.Context) (*bao.SealStatusResponse, error) {
	return c.client.Sys().SealStatusWithContext(ctx)
}
