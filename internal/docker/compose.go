package docker

import (
	"context"
	"os"
	"os/exec"
)

func (c *client) Compose(ctx context.Context, args ...string) error {

	cmdArgs := []string{
		"compose",
		"-f",
		c.cfg.Compose.File,
		"-p",
		c.cfg.Docker.Project,
	}

	cmdArgs = append(cmdArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
