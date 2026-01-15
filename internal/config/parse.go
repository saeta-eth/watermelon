package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// ParseFile reads and parses a .watermelon.toml file
func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	return Parse(data)
}

// Parse parses TOML config bytes
func Parse(data []byte) (*Config, error) {
	cfg := NewConfig()
	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return cfg, nil
}
