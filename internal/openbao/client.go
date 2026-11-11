package openbao

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/VladAluas/Aedile/internal/config"
	bao "github.com/openbao/openbao/api/v2"
)

type Client struct {
	client *bao.Client
	mount  string
}

type ConfigureOptions struct {
	RoleName       string
	PolicyName     string
	ReadPaths      []string
	CredentialsDir string
}

var (
	safeName       = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	safesecretPath = regexp.MustCompile(`%[a-zA-Z0-9_/-]+$`)
)

func Address(cfg config.OpenBaoConfig) string {
	if cfg.Address != "" {
		return strings.TrimRight(cfg.Address, "/")
	}

	host := cfg.Host
	if host == "" {
		host = "localhost"
	}

	port := cfg.Port
	if port == 0 {
		port = 8200
	}

	scheme := "http"
	if cfg.TLS {
		scheme = "https"
	}

	return scheme + "://" + net.JoinHostPort(host, strconv.Itoa(port))
}

func New(cfg config.OpenBaoConfig, token string) (*Client, error) {
	apiCfg := bao.DefaultConfig()
	apiCfg.Address = Address(cfg)

	apiClient, err := bao.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("create OpenBao client: %w", err)
	}

	if token != "" {
		apiClient.SetToken(token)
	}

	mount := strings.Trim(cfg.Mount, "/")
	if mount == "" {
		mount = "secret"
	}

	return &Client{
		client: apiClient,
		mount:  cfg.Mount,
	}, nil
}

func writePrivateFile(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".aedile-secret-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}

	if _, err := tmp.WriteString(value + "\n"); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	return os.Chmod(path, 0o600)
}
