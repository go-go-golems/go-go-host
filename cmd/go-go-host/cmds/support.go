package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/spf13/cobra"
)

const defaultAPIURL = "http://127.0.0.1:8080"

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
			MiddlewaresFunc:   glazedcli.CobraCommandDefaultMiddlewares,
		}),
	)
}

func getJSON(apiURL, path string, out any) error {
	return requestJSON(http.MethodGet, apiURL, path, "", "", nil, out)
}

func getJSONAsDevUser(apiURL, path, devUser string, out any) error {
	return requestJSON(http.MethodGet, apiURL, path, devUser, "", nil, out)
}

func getJSONWithAuth(apiURL, path, devUser, bearerToken string, out any) error {
	return requestJSON(http.MethodGet, apiURL, path, devUser, bearerToken, nil, out)
}

func postJSONAsDevUser(apiURL, path, devUser string, in, out any) error {
	return requestJSON(http.MethodPost, apiURL, path, devUser, "", in, out)
}

func postJSONWithAuth(apiURL, path, devUser, bearerToken string, in, out any) error {
	return requestJSON(http.MethodPost, apiURL, path, devUser, bearerToken, in, out)
}

func requestJSON(method, apiURL, path, devUser, bearerToken string, in, out any) error {
	base := strings.TrimRight(apiURL, "/")
	var body *bytes.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader(nil)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(method, base+path, body)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if devUser != "" {
		req.Header.Set("X-Go-Go-Host-User", devUser)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s%s: %w", method, base, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s %s%s: unexpected status %s", method, base, path, resp.Status)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
