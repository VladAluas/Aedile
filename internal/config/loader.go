// Package config helps translate the yaml file to docker
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	var cfg Config

	data, err := os.ReadFile("flow.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
