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

			return platform.
				New(cfg).
				Migrate(context.Background())
		},
	}
}
