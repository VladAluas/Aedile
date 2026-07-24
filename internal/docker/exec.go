package docker

import (
	"context"
	"os"
	"os/exec"
)

func (c *client) Exec(ctx context.Context, container string, args ...string) error {
	command := append([]string{"exec", container}, args...)

	cmd := exec.CommandContext(ctx, "docker", command...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
