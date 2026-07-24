package commands

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run DataForge",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			return platform.
				New(cfg).
				Run(context.Background())
		},
	}
}
