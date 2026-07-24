package platform

import "context"

func (p *Platform) Clean(ctx context.Context, all bool) error {
	args := []string{"down"}

	if all {
		args = append(args, "-v")
	}

	return p.docker.Compose(ctx, args...)
}
