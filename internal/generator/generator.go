// Package generator is responsible with generating the arthefacts for the services
package generator

import (
	"embed"
	"text/template"
	"os"

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

	return tmpl.Execute(file, cfg)
}
