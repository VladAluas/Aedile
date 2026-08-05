package commands

import (
	"context"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewDBCommand() *cobra.Command {
	db := &cobra.Command{
		Use:   "db",
		Short: "Database commands",
	}

	db.AddCommand(NewMigrateCommand())
	db.AddCommand(NewSeedCommand())

	return db
}

func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run Liquibase migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}
			return p.Migrate(context.Background())
		},
	}
}

func NewSeedCommand() *cobra.Command {
	return &cobra.Command{
		Use: "seed",
		Short: "Run Liquibase seed",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}
			return p.Seed(context.Background())
		},
	}
}
