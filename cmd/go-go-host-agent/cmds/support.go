package cmds

import (
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
	base := strings.TrimRight(apiURL, "/")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(base + path)
	if err != nil {
		return fmt.Errorf("GET %s%s: %w", base, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s%s: unexpected status %s", base, path, resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
