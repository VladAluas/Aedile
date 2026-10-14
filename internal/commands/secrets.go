package commands

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewSecretsCommand() *cobra.Command {
	var seed bool

	cmd := &cobra.Command{
		Use:   "secrests",
		Short: "Manage OpenBao Secrets",

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}

			if seed {
				return p.SeedSecrets(context.Background())
			}

			return cmd.Help()
		},
	}

	cmd.Flags().BoolVar(
		&seed,
		"seed",
		false,
		"Seed secreeets into OpenBao",
		)

	return cmd
}
