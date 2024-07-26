// Package configs offers config-related functionalit.
package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DataConfig stores configs.
type DataConfig struct {
	ListenAddress string
	DSN           string
}

// Obfuscate returns a string representation of the configs, without the security-risky entries.
func (dc *DataConfig) Obfuscate() (string, error) {
	aux := *dc
	// aux.RiskyField = "xxx"

	asJSON, err := json.MarshalIndent(aux, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not marshal: %w", err)
	}

	return string(asJSON), nil
}

// ReadConfigs reads configs from a JSON file.
func ReadConfigs(filename string) (*DataConfig, error) {
	file, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	result := DataConfig{}

	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}

	return &result, nil
}
