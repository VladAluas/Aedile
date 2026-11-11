// Package generator is responsible with generating the arthefacts for the services
package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/VladAluas/Aedile/internal/config"
)

//go:embed templates/*
var templates embed.FS

func GenerateCompose(cfg *config.Config) error {
	tmpl, err := template.ParseFS(
		templates,
		"templates/docker-compose.yaml.tmpl",
		"templates/docker-compose/*.tmpl",
	)
	if err != nil {
		return err
	}

	file, err := os.Create("infrastructure/docker-compose.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, cfg); err != nil {
		return fmt.Errorf("generate compose: %w", err)
	}

	if cfg.OpenBao.Enabled {
		if err := copyOpenBaoConfig(); err != nil {
			return err
		}
	}

	return nil
}

func copyOpenBaoConfig() error {
	data, err := templates.ReadFile("templates/openbao/config.hcl")
	if err != nil {
		return fmt.Errorf("read OpenBao config: %w", err)
	}

	path := "infrastructure/openbao/config.hcl"

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create OpenBao directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write OpenBao config: %w", err)
	}

	return nil
}
