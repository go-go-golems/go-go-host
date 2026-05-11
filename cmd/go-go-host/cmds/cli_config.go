package cmds

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CLIConfig struct {
	APIURL      string `yaml:"apiUrl" json:"apiUrl"`
	DevUser     string `yaml:"devUser" json:"devUser"`
	BearerToken string `yaml:"bearerToken" json:"bearerToken"`
}

func defaultCLIConfigPath() (string, error) {
	if path := os.Getenv("GO_GO_HOST_CLI_CONFIG"); path != "" {
		return path, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "go-go-host", "config.yaml"), nil
}

func loadCLIConfig() (CLIConfig, error) {
	path, err := defaultCLIConfigPath()
	if err != nil {
		return CLIConfig{}, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return CLIConfig{}, nil
	}
	if err != nil {
		return CLIConfig{}, fmt.Errorf("read CLI config %s: %w", path, err)
	}
	var cfg CLIConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return CLIConfig{}, fmt.Errorf("parse CLI config %s: %w", path, err)
	}
	return cfg, nil
}

func saveCLIConfig(cfg CLIConfig) (string, error) {
	path, err := defaultCLIConfigPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", fmt.Errorf("create CLI config directory: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal CLI config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("write CLI config %s: %w", path, err)
	}
	return path, nil
}

func resolveCLISettings(apiURL, devUser, bearerToken string) (CLIConfig, error) {
	cfg, err := loadCLIConfig()
	if err != nil {
		return CLIConfig{}, err
	}
	if apiURL != "" && apiURL != defaultAPIURL {
		cfg.APIURL = apiURL
	}
	if cfg.APIURL == "" {
		cfg.APIURL = defaultAPIURL
	}
	if devUser != "" {
		cfg.DevUser = devUser
	}
	if bearerToken != "" {
		cfg.BearerToken = bearerToken
	}
	return cfg, nil
}
