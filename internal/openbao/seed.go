package openbao

import (
	"context"
	"fmt"
	"strings"
)

// Client Seed
func (c *Client) Seed(ctx context.Context, secrets map[string]map[string]string) error {
	kv := c.client.KVv2(c.mount)

	for path, values := range secrets {
		path = strings.Trim(path, "/")
		if !safesecretPath.MatchString(path) || path == "" || strings.Contains(path, "..") {
			return fmt.Errorf("invalid secret path %q", path)
		}
		if len(values) == 0 {
			return fmt.Errorf("secret %q has no values", path)
		}

		data := make(map[string]any, len(values))
		for key, value := range values {
			data[key] = value
		}

		if _, err := kv.Put(ctx, path, data); err != nil {
			return fmt.Errorf("seed secret %q: %w", path, err)
		}
	}

	return nil
}
