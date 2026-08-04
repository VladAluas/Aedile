package platform

import (
	"database/sql"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/db"
	"github.com/VladAluas/Aedile/internal/docker"
)

type Platform struct {
	cfg    *config.Config
	docker docker.Client
	db     *sql.DB
}

func (p *Platform) DB() (*sql.DB, error) {
	if p.db != nil {
		return p.db, nil
	}

	conn, err := db.Connect(p.cfg)
	if err != nil {
		return nil, err
	}

	p.db = conn
	return p.db, nil
}

func New(cfg *config.Config) (*Platform, error) {
	return &Platform{
			cfg:    cfg,
			docker: docker.New(cfg),
		},
		nil
}

func (p *Platform) Close() error {
	if p.db == nil {
		return nil
	}

	return p.db.Close()
}
