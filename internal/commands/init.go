package commands

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}
			return p.Init(context.Background())
		},
	}
}
