package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SecurityAction represents the result of a security check.
type SecurityAction int

const (
	ActionAllow SecurityAction = iota
	ActionBlock
	ActionAsk
)

// Config holds the application configuration.
type Config struct {
	Blocked []string `yaml:"blocked"`
	Allowed []string `yaml:"allowed"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Blocked: []string{
			"*.exe", "*.bat", "*.cmd", "*.ps1", "*.vbs", "*.js", "*.msi",
			"*.scr", "*.com", "*.pif", "*.reg", "*.wsf", "*.wsh",
		},
		Allowed: []string{},
	}
}

// ConfigPath returns the path to the config file.
func ConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "ope", "ope.yml"), nil
}

// LoadConfig reads the config from disk, or returns default config if not found.
func LoadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// SaveConfig writes the config to disk.
func SaveConfig(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// CheckSecurity determines the security action for a given path.
func (c *Config) CheckSecurity(path string) SecurityAction {
	base := strings.ToLower(filepath.Base(path))

	for _, pattern := range c.Blocked {
		if matched, _ := filepath.Match(strings.ToLower(pattern), base); matched {
			return ActionBlock
		}
	}

	for _, pattern := range c.Allowed {
		if matched, _ := filepath.Match(strings.ToLower(pattern), base); matched {
			return ActionAllow
		}
	}

	// Directories are allowed by default
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return ActionAllow
	}

	return ActionAsk
}
