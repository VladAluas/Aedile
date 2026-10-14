package openbao

import (
	"fmt"

	"github.com/VladAluas/Aedile/internal/config"
	bao "github.com/openbao/openbao/api/v2"
)

type Client struct {
	client *bao.Client
	mount  string
}

func New(cfg config.OpenBaoConfig, token string) (*Client, error) {
	baoCfg := bao.DefaultConfig()
	baoCfg.Address = fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	client, err := bao.NewClient(baoCfg)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	return &Client{
		client: client,
		mount: cfg.Mount,
	}, nil
}
