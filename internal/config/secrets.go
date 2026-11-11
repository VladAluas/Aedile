package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type SecretValues map[string]map[string]string

func LoadSecretsFile(path string) (SecretValues, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read secrets file: %w", err)
	}

	var secrets SecretValues
	if err := yaml.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("parse secrets file: %w", err)
	}

	if len(secrets) == 0 {
		return nil, fmt.Errorf("secrets file %q contains no secrets", path)
	}

	for path, values := range secrets {
		if strings.TrimSpace(path) == "" || len(values) == 0 {
			return nil, fmt.Errorf("secrets path %q has no values", path)
		}
	}

	return secrets, nil
}
