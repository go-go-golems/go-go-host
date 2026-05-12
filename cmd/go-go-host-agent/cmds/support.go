package cmds

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/go-go-host/internal/control"
	"github.com/spf13/cobra"
)

const defaultAPIURL = "http://127.0.0.1:8080"

type AgentConfig struct {
	APIURL     string `json:"apiUrl"`
	AgentID    string `json:"agentId"`
	KeyID      string `json:"keyId"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	OrgID      string `json:"orgId"`
	SiteID     string `json:"siteId"`
}

func standardSections() (schema.Section, schema.Section, error) {
	glazedSection, err := settings.NewGlazedSchema()
	if err != nil {
		return nil, nil, err
	}
	commandSettingsSection, err := glazedcli.NewCommandSettingsSection()
	if err != nil {
		return nil, nil, err
	}
	return glazedSection, commandSettingsSection, nil
}

func BuildGlazedCobraCommand(command glazedcmds.Command) (*cobra.Command, error) {
	return glazedcli.BuildCobraCommandFromCommand(command,
		glazedcli.WithParserConfig(glazedcli.CobraParserConfig{
			ShortHelpSections: []string{schema.DefaultSlug},
			AppName:           "GO_GO_HOST_AGENT",
		}),
	)
}

func getJSON(apiURL, path string, out any) error {
	base := strings.TrimRight(apiURL, "/")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(base + path)
	if err != nil {
		return fmt.Errorf("GET %s%s: %w", base, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s%s: unexpected status %s: %s", base, path, resp.Status, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func postJSON(apiURL, path string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	base := strings.TrimRight(apiURL, "/")
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Post(base+path, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST %s%s: unexpected status %s: %s", base, path, resp.Status, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func signedPostJSON(apiURL, path string, body any, cfg AgentConfig, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	base := strings.TrimRight(apiURL, "/")
	req, err := http.NewRequest(http.MethodPost, base+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err := signRequest(req, payload, cfg); err != nil {
		return err
	}
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST %s: unexpected status %s: %s", path, resp.Status, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func uploadBundle(apiURL, runID, uploadToken, bundlePath string, out any) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("bundle", filepath.Base(bundlePath))
	if err != nil {
		return err
	}
	f, err := os.Open(bundlePath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, f); err != nil {
		_ = f.Close()
		return err
	}
	_ = f.Close()
	_ = writer.Close()
	base := strings.TrimRight(apiURL, "/")
	req, err := http.NewRequest(http.MethodPost, base+"/api/v1/agent/deploy-runs/"+url.PathEscape(runID)+"/upload", &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Go-Go-Upload-Token", uploadToken)
	resp, err := (&http.Client{Timeout: 2 * time.Minute}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload: unexpected status %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func signRequest(req *http.Request, body []byte, cfg AgentConfig) error {
	privBytes, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		return fmt.Errorf("decode private key: %w", err)
	}
	if len(privBytes) != ed25519.PrivateKeySize {
		return fmt.Errorf("private key must decode to %d bytes", ed25519.PrivateKeySize)
	}
	timestamp := time.Now().UTC().Format(time.RFC3339)
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return err
	}
	nonce := base64.RawURLEncoding.EncodeToString(nonceBytes)
	canonical := control.AgentCanonicalString(req.Method, req.URL.RequestURI(), hashBody(body), timestamp, nonce)
	sig := ed25519.Sign(ed25519.PrivateKey(privBytes), []byte(canonical))
	req.Header.Set("X-Go-Go-Agent-ID", cfg.AgentID)
	req.Header.Set("X-Go-Go-Agent-Key-ID", cfg.KeyID)
	req.Header.Set("X-Go-Go-Agent-Timestamp", timestamp)
	req.Header.Set("X-Go-Go-Agent-Nonce", nonce)
	req.Header.Set("X-Go-Go-Agent-Signature", base64.StdEncoding.EncodeToString(sig))
	return nil
}

func hashBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func defaultConfigPath() string {
	if v := os.Getenv("GO_GO_HOST_AGENT_CONFIG"); v != "" {
		return v
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return ".go-go-host-agent.json"
	}
	return filepath.Join(dir, "go-go-host-agent", "config.json")
}

func loadConfig(path string) (AgentConfig, error) {
	if path == "" {
		path = defaultConfigPath()
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return AgentConfig{}, err
	}
	var cfg AgentConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return AgentConfig{}, err
	}
	return cfg, nil
}

func saveConfig(path string, cfg AgentConfig) error {
	if path == "" {
		path = defaultConfigPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(path, b, 0o600)
}

func generateKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv), nil
}
