package openbao

import (
	"context"
	"fmt"
)

func (c *Client) Seed(
	ctx context.Context,
	secrests map[string]map[string]any,
) error {
	kv := c.client.KVv2(c.mount)

	for path, values := range secrests {
		if _, err := kv.Put(ctx, path, values); err != nil {
			return fmt.Errorf("failed to seed secret %q: %w",
				path,
				err,
			)
		}
	}

	return nil
}
