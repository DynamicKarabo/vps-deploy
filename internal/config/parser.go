package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parse reads and parses a deploy.yaml file into a Config struct.
func Parse(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Inject service names into each Service struct.
	for name, svc := range cfg.Services {
		if svc != nil {
			svc.Name = name
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	for name, svc := range c.Services {
		if svc == nil {
			return fmt.Errorf("service %q is empty", name)
		}
		if svc.Host == "" {
			return fmt.Errorf("service %q: host is required", name)
		}
		if svc.User == "" {
			return fmt.Errorf("service %q: user is required", name)
		}
	}
	return nil
}
