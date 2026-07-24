// Package platform orchestrates workflows
package platform

import "context"

func (p *Platform) Init(ctx context.Context) error {
	return p.docker.Compose(ctx, "up", "-d")
}
