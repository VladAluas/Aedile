// Package metadata helps load the metadata config from /configs/platform.json to the metadata db
package metadata

import (
	"encoding/json"
	"os"
)

const configFile = "configs/platform.json"

func Load() (*Config, error) {
	var cfg Config

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
