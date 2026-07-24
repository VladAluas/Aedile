// Package docker is intended to help execute the docker commands
package docker

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
)

type Client interface {
	Compose(ctx context.Context, args ...string) error
	Exec(ctx context.Context, container string, args ...string) error
}

func New(cfg *config.Config) Client {
	return &client{cfg: cfg}
}

type client struct {
	cfg *config.Config
}
