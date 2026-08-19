package platform

import (
	"context"
)

func (p *Platform) Migrate(ctx context.Context) error {
	return p.docker.Compose(
		ctx,
		"run",
		"--rm",
		"liquibase",
		"--defaults-file=/liquibase/liquibase.properties",
		"update",
	)
}
