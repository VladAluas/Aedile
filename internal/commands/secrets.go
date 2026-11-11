package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/VladAluas/Aedile/internal/config"
	"github.com/VladAluas/Aedile/internal/openbao"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

func NewSecretsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage OpenBao secrets",
	}

	cmd.AddCommand(
		NewSecretsStatusCommand(),
		NewSecretsInitCommand(),
		NewSecretsUnsealCommand(),
		NewSecretsConfigureCommand(),
		NewSecretsSeedCommand(),
	)

	return cmd
}

func secretsClient(authenticated bool) (*config.Config, *openbao.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	address := openbao.Address(cfg.OpenBao)
	production := strings.EqualFold(cfg.Project.Environment, "prod") || strings.EqualFold(cfg.Project.Environment, "production")

	if production && !strings.HasPrefix(strings.ToLower(address), "https://") {
		return nil, nil, fmt.Errorf(
			"production OpenBao requires a HTTPS address in openbao.address",
		)
	}

	token := ""
	if authenticated {
		token = os.Getenv("OPENBAO_TOKEN")
		if token == "" {
			return nil, nil, fmt.Errorf(
				"OPENBAO_TOKEN is required for this operation",
			)
		}
	}

	client, err := openbao.New(cfg.OpenBao, token)
	if err != nil {
		return nil, nil, err
	}

	return cfg, client, nil
}

func NewSecretsStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show OpenBao initialization and seal status",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := secretsClient(false)
			if err != nil {
				return err
			}

			status, err := client.Status(cmd.Context())
			if err != nil {
				return fmt.Errorf("query OpenBao status: %w", err)
			}

			fmt.Printf("Initialized: %t\n", status.Initialized)
			fmt.Printf("Sealed:      %t\n", status.Sealed)
			fmt.Printf("Version:     %s\n", status.Version)

			if status.Initialized {
				fmt.Printf("Unseal:      %d/%d\n", status.Progress, status.T)
				fmt.Printf("Key shares:  %d\n", status.N)
				fmt.Printf("Threshold:   %d\n", status.T)
			}

			return nil
		},
	}
}

func NewSecretsInitCommand() *cobra.Command {
	var shares, threshold int
	var pgpKeys []string
	var rootTokenPGPKey string
	var yes bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize an uninitialized OpenBao server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := secretsClient(false)
			if err != nil {
				return err
			}

			production := strings.EqualFold(cfg.Project.Environment, "prod") || strings.EqualFold(cfg.Project.Environment, "production")
			if production && (len(pgpKeys) != shares || rootTokenPGPKey == "") {
				return fmt.Errorf(
					"production initialization requires %d --pgp-keys and --root-token-pgp-keyl",
					shares,
				)
			}

			if !yes {
				fmt.Fprintln(cmd.ErrOrStderr(),
					"WARNING: intitialization is a ont-time operation.")
				fmt.Fprintln(cmd.ErrOrStderr(),
					"Protect the generated unseal keys and root token separately.")
				fmt.Fprint(cmd.ErrOrStderr(), "Initialize OpenBao now? [y/N] ")

				reader := bufio.NewReader(cmd.InOrStdin())
				answer, _ := reader.ReadString('\n')
				if !strings.EqualFold(strings.TrimSpace(answer), "y") && !strings.EqualFold(strings.TrimSpace(answer), "yes") {
					return fmt.Errorf("initialization cancelled")
				}
			}

			result, err := client.Initialize(
				cmd.Context(),
				shares,
				threshold,
				pgpKeys,
				rootTokenPGPKey,
			)
			if err != nil {
				return err
			}

			if !production && len(pgpKeys) == 0 {
				fmt.Fprintln(cmd.ErrOrStderr(),
					"WARNING: the following output contains plaintext key material. Store it securely.",
				)
			}

			output := struct {
				Keys       []string `json:"keys,omitempty"`
				KeysBase64 []string `json:"keys_base64,omitempty"`
				RootToken  string   `json:"root_token,omitempty"`
			}{
				Keys:       result.Keys,
				KeysBase64: result.KeysB64,
				RootToken:  result.RootToken,
			}

			data, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStderr(), string(data))
			fmt.Fprintln(cmd.OutOrStderr(),
				"Save the initialization material securely. It cannot be retrieved by rerunning init.")

			return nil
		},
	}

	cmd.Flags().IntVar(&shares, "key-shares", 5, "Number of unseal key shares")
	cmd.Flags().IntVar(&threshold, "key-threshold", 3, "Number of shares required to unseal")
	cmd.Flags().StringSliceVar(&pgpKeys, "pgp-keys", nil, "PGP public-key files, one per key share")
	cmd.Flags().StringVar(&rootTokenPGPKey, "root-token-pgp-key", "", "PGP public key file for encrypting the root token")
	cmd.Flags().BoolVar(&yes, "yes", false, "Skip interactive confirmation")

	return cmd
}

func NewSecretsUnsealCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "unseal",
		Short: "Unseal OpenBao by prompting for key shares",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := secretsClient(false)
			if err != nil {
				return err
			}

			status, err := client.Status(cmd.Context())
			if err != nil {
				return err
			}
			if !status.Initialized {
				return fmt.Errorf("OpenBao is not initialized; run 'aedile secrets init' first")
			}
			if !status.Sealed {
				fmt.Fprintln(cmd.OutOrStdout(), "OpenBao is already unsealed.")
				return nil
			}

			maxAttempts := status.N + status.T
			for attempt := 0; status.Sealed && attempt < maxAttempts; attempt++ {
				fmt.Fprintf(
					cmd.ErrOrStderr(),
					"Enter unseal key (progress %d/%d): ",
					status.Progress,
					status.T,
				)

				key, err := term.ReadPassword(os.Stdin.Fd())
				fmt.Fprintln(cmd.ErrOrStderr())
				if err != nil {
					return fmt.Errorf("read unseal key: %w", err)
				}

				status, err = client.Unseal(cmd.Context(), strings.TrimSpace(string(key)))
				for i := range key {
					key[i] = 0
				}
				if err != nil {
					return fmt.Errorf("submit unseal key: %w", err)
				}
			}

			if status.Sealed {
				return fmt.Errorf("OpenBao remains sealed; progress %d/%d", status.Progress, status.T)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OpenBao is unsealed")
			return nil
		},
	}
}

func NewSecretsConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Configure KV v2 Sisyphus AppRole",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, client, err := secretsClient(false)
			if err != nil {
				return err
			}

			status, err := client.Status(cmd.Context())
			if err != nil {
				return err
			}
			if !status.Initialized || status.Sealed {
				return fmt.Errorf("OpenBao must be initialized and unsealed before configuring")
			}

			opts := openbao.ConfigureOptions{
				RoleName:   cfg.OpenBao.AppRole,
				PolicyName: cfg.OpenBao.Policy,
				ReadPaths:  cfg.OpenBao.ReadPaths,
			}
			if err := client.Configure(cmd.Context(), opts); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OpenBao KV v2 and Approle configured")
			fmt.Fprintln(cmd.OutOrStdout(), "AppRole credentials are in infrastructure/secrets/. Keep these files private.")

			return nil
		},
	}
}

func NewSecretsSeedCommand() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed OpenBao KV v2 from aedile.secrets.yaml",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, client, err := secretsClient(true)
			if err != nil {
				return err
			}

			status, err := client.Status(cmd.Context())
			if err != nil {
				return err
			}
			if !status.Initialized || status.Sealed {
				return fmt.Errorf("OpenBao must be initialized and unsealed before seeding")
			}

			secrets, err := config.LoadSecretsFile(file)
			if err != nil {
				return err
			}
			if err := client.Seed(cmd.Context(), secrets); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(),
				"Seeded %d secret paths into OpenBao KV v2. \n",
				len(secrets),
			)

			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "aedile.secrets.yaml", "Secrets YAML file")
	return cmd
}
