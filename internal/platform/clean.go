package platform

import "context"

func (p *Platform) Clean(ctx context.Context, all bool, images bool) error {
	args := []string{"down"}

	if all {
		args = append(args, "-v")
	}

	if images {
		args = append(args, "-rmi", "all")
	}

	return p.docker.Compose(ctx, args...)
}
