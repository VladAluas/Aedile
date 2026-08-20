package commands

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewCleanCommand() *cobra.Command {
	var all 	 bool
	var images bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Stop the platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			p, err := platform.New(cfg)
			if err != nil {
				return err
			}

			return p.Clean(context.Background(), all, images)
		},
	}

	cmd.Flags().BoolVar(
		&all,
		"all",
		false,
		"Remove volumes",
	)

	cmd.Flags().BoolVar(
		&images,
		"images",
		false,
		"Remove images",
	)

	return cmd
}
