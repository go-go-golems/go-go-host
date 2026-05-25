package cmds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type CLIConfig struct {
	APIURL      string          `yaml:"apiUrl" json:"apiUrl"`
	DevUser     string          `yaml:"devUser,omitempty" json:"devUser,omitempty"`
	BearerToken string          `yaml:"bearerToken,omitempty" json:"bearerToken,omitempty"`
	OIDC        *CLIOIDCSession `yaml:"oidc,omitempty" json:"oidc,omitempty"`
}

type CLIOIDCSession struct {
	Issuer       string    `yaml:"issuer" json:"issuer"`
	ClientID     string    `yaml:"clientId" json:"clientId"`
	Scopes       []string  `yaml:"scopes,omitempty" json:"scopes,omitempty"`
	AccessToken  string    `yaml:"accessToken" json:"accessToken"`
	IDToken      string    `yaml:"idToken,omitempty" json:"idToken,omitempty"`
	RefreshToken string    `yaml:"refreshToken,omitempty" json:"refreshToken,omitempty"`
	TokenType    string    `yaml:"tokenType,omitempty" json:"tokenType,omitempty"`
	ExpiresAt    time.Time `yaml:"expiresAt,omitempty" json:"expiresAt,omitempty"`
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
	// #nosec G117 -- CLI config intentionally persists OAuth credentials with 0600 permissions.
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
	return resolveCLISettingsContext(context.Background(), apiURL, devUser, bearerToken)
}

func resolveCLISettingsContext(ctx context.Context, apiURL, devUser, bearerToken string) (CLIConfig, error) {
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
	if cfg.DevUser != "" || cfg.BearerToken != "" || cfg.OIDC == nil || cfg.OIDC.AccessToken == "" {
		return cfg, nil
	}
	if tokenExpiresSoon(cfg.OIDC.ExpiresAt) && cfg.OIDC.RefreshToken != "" {
		refreshed, err := refreshOIDCToken(ctx, cfg.OIDC)
		if err != nil {
			return CLIConfig{}, err
		}
		cfg.OIDC = refreshed
		if _, err := saveCLIConfig(cfg); err != nil {
			return CLIConfig{}, err
		}
	}
	cfg.BearerToken = cfg.OIDC.AccessToken
	return cfg, nil
}

func tokenExpiresSoon(expiresAt time.Time) bool {
	if expiresAt.IsZero() {
		return false
	}
	return time.Until(expiresAt) < 60*time.Second
}
