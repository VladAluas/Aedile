// Package commands contains the CLI commands for the infrastructure
package commands

import "github.com/spf13/cobra"

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "aedile",
		Short: "Developer CLI for the DataForge platform",
		Long: `Flow manages the local DataForge platform.
					 It provides commands to manage infrastructure,
					 database migrations, and ETL execution.
					`,
	}

	root.AddCommand(NewInitCommand())
	root.AddCommand(NewCleanCommand())
	root.AddCommand(NewRunCommand())
	root.AddCommand(NewDBCommand())
	root.AddCommand(NewServicesCommand())

	return root
}
