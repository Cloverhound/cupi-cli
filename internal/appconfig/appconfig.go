package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CredentialConfig stores username for a credential type
type CredentialConfig struct {
	Username string `json:"username"`
}

// ServerConfig holds CUC server configuration
type ServerConfig struct {
	Host        string                      `json:"host"`
	Port        int                         `json:"port"`
	Version     string                      `json:"version"`
	Credentials map[string]CredentialConfig `json:"credentials"`
}

// Config is the root configuration structure
type Config struct {
	DefaultServer string                  `json:"defaultServer"`
	Servers       map[string]ServerConfig `json:"servers"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.cupi-cli/config.json"
	}
	return filepath.Join(home, ".cupi-cli", "config.json")
}

// LoadConfig loads the configuration from ~/.cupi-cli/config.json
// Returns empty Config if file doesn't exist
func LoadConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				DefaultServer: "",
				Servers:       make(map[string]ServerConfig),
			}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Servers == nil {
		cfg.Servers = make(map[string]ServerConfig)
	}

	return &cfg, nil
}

// SaveConfig persists configuration to ~/.cupi-cli/config.json
// Creates the directory if it doesn't exist
func SaveConfig(cfg *Config) error {
	path := configPath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetServer retrieves a server configuration by name
func GetServer(cfg *Config, name string) (*ServerConfig, error) {
	server, ok := cfg.Servers[name]
	if !ok {
		return nil, fmt.Errorf("server '%s' not found", name)
	}
	return &server, nil
}

// SetDefaultServer sets the default server name
func SetDefaultServer(name string) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	if _, ok := cfg.Servers[name]; !ok {
		return fmt.Errorf("server '%s' not found", name)
	}

	cfg.DefaultServer = name
	return SaveConfig(cfg)
}
