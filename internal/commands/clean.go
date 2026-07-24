package commands

import (
	"context"

	"github.com/VladAluas/flow/internal/config"
	"github.com/VladAluas/flow/internal/platform"
	"github.com/spf13/cobra"
)

func NewCleanCommand() *cobra.Command {

	var all bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Stop the platform",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			return platform.
				New(cfg).
				Clean(context.Background(), all)
		},
	}

	cmd.Flags().BoolVar(
		&all,
		"all",
		false,
		"Remove volumes",
	)

	return cmd
}
