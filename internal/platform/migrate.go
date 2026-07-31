package platform

import (
	"context"
)

func (p *Platform) Migrate(ctx context.Context) error {
	return p.docker.Compose(ctx, "--profile", "tools", "run", "--rm", "liquibase")
}
