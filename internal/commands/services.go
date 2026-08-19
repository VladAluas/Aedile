package commands

import (
	"context"
	"fmt"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/platform"
	"github.com/spf13/cobra"
)

func NewServicesCommand() *cobra.Command {
	services := &cobra.Command{
		Use:   "services",
		Short: "Services commands",
	}

	services.AddCommand(ListCommand())
	services.AddCommand(DescribeCommand())
	// Need to implement at a later date
	//
	// services.AddCommand(ValidateCommand())

	return services
}

func DescribeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "describe <template>",
		Short: "Describe the configuration needed for a specific service",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("template name is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}

			return p.DescribeService(context.Background(), args[0])
		},
	}
}

func ListCommand() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		Short: "List all the services needed for inital configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			p, err := platform.New(cfg)
			if err != nil {
				return err
			}

			return p.ListServices(context.Background())
		},
	}
}

// Need to implement at a later date
//
// func ValidateCommand() *cobra.Command {
// 	return &cobra.Command{
// 		Use: "validate",
// 		Short: "Validate the aedile.yaml file configuration",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cfg, err := config.Load()
// 			if err != nil {
// 				return err
// 			}
//
// 			p, err := platform.New(cfg)
// 			if err != nil {
// 				return err
// 			}
//
// 			return p.ValidateServices(context.Background())
// 		},
// 	}
// }
