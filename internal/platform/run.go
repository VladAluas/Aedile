package platform

import "context"

func (p *Platform) Run(ctx context.Context, args ...string) error {
	cmdArgs := []string{
		"run",
		"--rm",
		"etl",
	}

	cmdArgs = append(cmdArgs, args...)

	return p.docker.Compose(ctx, cmdArgs...)
}
