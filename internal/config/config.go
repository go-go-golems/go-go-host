package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config contains daemon settings used by the phase-0 skeleton. Later phases
// will extend this with database, OIDC, runtime, and quota settings.
type Config struct {
	ListenAddr                string        `json:"listenAddr" yaml:"listenAddr"`
	PublicBaseURL             string        `json:"publicBaseUrl" yaml:"publicBaseUrl"`
	BaseDomain                string        `json:"baseDomain" yaml:"baseDomain"`
	DataDir                   string        `json:"dataDir" yaml:"dataDir"`
	ControlDBDSN              string        `json:"controlDbDsn" yaml:"controlDbDsn"`
	OIDCIssuer                string        `json:"oidcIssuer" yaml:"oidcIssuer"`
	OIDCClientID              string        `json:"oidcClientId" yaml:"oidcClientId"`
	OIDCDeviceClientID        string        `json:"oidcDeviceClientId" yaml:"oidcDeviceClientId"`
	OIDCAcceptedClientIDs     []string      `json:"oidcAcceptedClientIds" yaml:"oidcAcceptedClientIds"`
	OIDCScopes                []string      `json:"oidcScopes" yaml:"oidcScopes"`
	OIDCRedirectPath          string        `json:"oidcRedirectPath" yaml:"oidcRedirectPath"`
	OIDCLogoutRedirectPath    string        `json:"oidcLogoutRedirectPath" yaml:"oidcLogoutRedirectPath"`
	DevAuth                   bool          `json:"devAuth" yaml:"devAuth"`
	DevPlatformAdminSubjects  []string      `json:"devPlatformAdminSubjects" yaml:"devPlatformAdminSubjects"`
	PlatformAdminOIDCSubjects []string      `json:"platformAdminOIDCSubjects" yaml:"platformAdminOIDCSubjects"`
	PlatformAdminEmails       []string      `json:"platformAdminEmails" yaml:"platformAdminEmails"`
	PlatformAdminOIDCRoles    []string      `json:"platformAdminOIDCRoles" yaml:"platformAdminOIDCRoles"`
	LogLevel                  string        `json:"logLevel" yaml:"logLevel"`
	ReadTimeout               time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout              time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	ShutdownTimeout           time.Duration `json:"shutdownTimeout" yaml:"shutdownTimeout"`
}

// Default returns a local-development configuration.
func Default() Config {
	return Config{
		ListenAddr:               "127.0.0.1:8080",
		PublicBaseURL:            "http://127.0.0.1:8080",
		BaseDomain:               "localhost",
		DataDir:                  "./data",
		ControlDBDSN:             "postgres://go_go_host:go_go_host_dev@127.0.0.1:55432/go_go_host?sslmode=disable",
		OIDCDeviceClientID:       "go-go-host-cli",
		OIDCScopes:               []string{"openid", "profile", "email"},
		OIDCRedirectPath:         "/app/auth/callback",
		OIDCLogoutRedirectPath:   "/app",
		DevAuth:                  true,
		DevPlatformAdminSubjects: []string{"dev-user"},
		LogLevel:                 "info",
		ReadTimeout:              15 * time.Second,
		WriteTimeout:             30 * time.Second,
		ShutdownTimeout:          5 * time.Second,
	}
}

// Load reads a YAML or JSON config file and overlays it on top of defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	if strings.TrimSpace(path) == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	data = []byte(os.ExpandEnv(string(data)))
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse json config %s: %w", path, err)
		}
	default:
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse yaml config %s: %w", path, err)
		}
	}
	cfg.ApplyDefaults()
	return cfg, nil
}

// ApplyDefaults fills fields omitted by a config file.
func (c *Config) ApplyDefaults() {
	defaults := Default()
	if c.ListenAddr == "" {
		c.ListenAddr = defaults.ListenAddr
	}
	if c.PublicBaseURL == "" {
		c.PublicBaseURL = defaults.PublicBaseURL
	}
	if c.BaseDomain == "" {
		c.BaseDomain = defaults.BaseDomain
	}
	if c.DataDir == "" {
		c.DataDir = defaults.DataDir
	}
	if c.ControlDBDSN == "" {
		c.ControlDBDSN = defaults.ControlDBDSN
	}
	if c.OIDCDeviceClientID == "" {
		c.OIDCDeviceClientID = defaults.OIDCDeviceClientID
	}
	if len(c.OIDCScopes) == 0 {
		c.OIDCScopes = defaults.OIDCScopes
	}
	if len(c.OIDCAcceptedClientIDs) == 0 && c.OIDCClientID != "" {
		c.OIDCAcceptedClientIDs = []string{c.OIDCClientID}
		if c.OIDCDeviceClientID != "" && c.OIDCDeviceClientID != c.OIDCClientID {
			c.OIDCAcceptedClientIDs = append(c.OIDCAcceptedClientIDs, c.OIDCDeviceClientID)
		}
	}
	if c.OIDCRedirectPath == "" {
		c.OIDCRedirectPath = defaults.OIDCRedirectPath
	}
	if c.OIDCLogoutRedirectPath == "" {
		c.OIDCLogoutRedirectPath = defaults.OIDCLogoutRedirectPath
	}
	if c.LogLevel == "" {
		c.LogLevel = defaults.LogLevel
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaults.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaults.WriteTimeout
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = defaults.ShutdownTimeout
	}
}
