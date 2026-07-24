package platform

import (
	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/docker"
)

type Platform struct {
	cfg    *config.Config
	docker docker.Client
}

func New(cfg *config.Config) *Platform {
	return &Platform{
		cfg:    cfg,
		docker: docker.New(cfg),
	}
}
