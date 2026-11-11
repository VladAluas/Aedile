package openbao

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	bao "github.com/openbao/openbao/api/v2"
)

// Client Configure
func (c *Client) Configure(ctx context.Context, opts ConfigureOptions) error {
	if opts.RoleName == "" {
		opts.RoleName = "sisyphus"
	}

	if opts.PolicyName == "" {
		opts.PolicyName = "sisyphus-read"
	}

	if opts.CredentialsDir == "" {
		opts.CredentialsDir = "infrastructure/secrets"
	}

	if len(opts.ReadPaths) == 0 {
		opts.ReadPaths = []string{"metadata", "storage", "etl"}
	}

	if !safeName.MatchString(opts.RoleName) || !safeName.MatchString(opts.PolicyName) {
		return fmt.Errorf("invalid AppName or policy name")
	}

	// Ensure the configured KV exists and it's version 2
	mounts, err := c.client.Sys().ListMountsWithContext(ctx)
	if err != nil {
		return fmt.Errorf("list secret mounts: %w", err)
	}

	mountPath := c.mount + "/"
	if existing, ok := mounts[mountPath]; ok {
		if existing.Type != "kv" || existing.Options["version"] != "2" {
			return fmt.Errorf("mount %q exists, but is not KV v2; refusing to modify it", mountPath)
		}
	} else {
		err = c.client.Sys().MountWithContext(ctx, mountPath, &bao.MountInput{
			Type:        "kv",
			Description: "Aedile application secrets",
			Options:     map[string]string{"version": "2"},
		})
		if err != nil {
			return fmt.Errorf("enable KV v2: %w", err)
		}
	}

	// Enable AppRole only if it is not already enabled.
	auths, err := c.client.Sys().ListAuthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("list authentication methods: %w", err)
	}

	if existing, ok := auths["approle/"]; ok {
		if existing.Type != "approle" {
			return fmt.Errorf("auth/approle is ocupied by another auth method")
		}
	} else {
		if err := c.client.Sys().EnableAuthWithOptionsWithContext(
			ctx,
			"approle",
			&bao.EnableAuthOptions{
				Type:        "approle",
				Description: "Aedile workload authentication",
			},
		); err != nil {
			return fmt.Errorf("enable AppRole: %w", err)
		}
	}

	// The Sisyphus policy grants read access only to selected secret paths
	var rules strings.Builder

	for _, path := range opts.ReadPaths {
		path = strings.Trim(path, "/")
		if !safesecretPath.MatchString(path) ||
			path == "" ||
			strings.Contains(path, "..") {
			return fmt.Errorf("invalid OpenBao read path: %q", err)
		}
		fmt.Fprintf(
			&rules,
			"path %q {\n  capabilities = [\"read\"]\n}\n\n",
			c.mount+"/data/"+path,
		)
	}

	if err := c.client.Sys().PutPolicyWithContext(
		ctx,
		opts.PolicyName,
		rules.String(),
	); err != nil {
		return fmt.Errorf("write AppRole policy: %w", err)
	}

	// Reapply the role definition on every run.
	rolePath := "auth/approle/role/" + opts.RoleName
	_, err = c.client.Logical().WriteWithContext(ctx, rolePath, map[string]any{
		"token_policies":     []string{opts.PolicyName},
		"token_ttl":          "1h",
		"token_max_ttl":      "4h",
		"token_num_uses":     0,
		"bind_secret_id":     true,
		"secret_id_ttl":      "720h",
		"secret_id_num_uses": 0,
	})
	if err != nil {
		return fmt.Errorf("configure AppRole: %w", err)
	}

	// Retrieve the RoleID
	role, err := c.client.Logical().ReadWithContext(ctx, rolePath+"/role-id")
	if err != nil {
		return fmt.Errorf("read AppRole ID: %w", err)
	}
	if role == nil || role.Data == nil {
		return fmt.Errorf("OpenBao returned no RoleID")
	}

	roleID, ok := role.Data["role_id"].(string)
	if !ok || roleID == "" {
		return fmt.Errorf("OpenBao returned and invalid RoleID")
	}

	if err := writePrivateFile(
		filepath.Join(opts.CredentialsDir, "openbao_role_id"),
		roleID,
	); err != nil {
		return fmt.Errorf("save AppRole ID: %w", err)
	}

	// Create the SecretID only if a local credential file does not already
	// exist. This make conf repeatable without rotating credentials on every
	// invocation.
	secretIDPath := filepath.Join(opts.CredentialsDir, "openbao_secret_id")
	secretIDBytes, err := os.ReadFile(secretIDPath)

	switch {
	case err == nil:
		if strings.TrimSpace(string(secretIDBytes)) == "" {
			return fmt.Errorf("%s exists but is empty. Remove it to issue a new secretID", secretIDPath)
		}
	case errors.Is(err, os.ErrNotExist):
		secret, err := c.client.Logical().WriteWithContext(
			ctx,
			rolePath+"/secret-id",
			map[string]any{},
		)
		if err != nil {
			return fmt.Errorf("create AppRole SecretID: %w", err)
		}
		if secret == nil || secret.Data == nil {
			return fmt.Errorf("OpenBao returned no SecretID")
		}

		secretID, ok := secret.Data["secret_id"].(string)
		if !ok || secretID == "" {
			return fmt.Errorf("OpenBao returned an invalid SecretID")
		}

		if err := writePrivateFile(secretIDPath, secretID); err != nil {
			return fmt.Errorf("save AppRole SecretID: %w", err)
		}

	default:
		return fmt.Errorf("read AppRole SecretID file: %w", err)
	}

	return nil
}
