package cmds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewMaintenanceCobraCommand() (*cobra.Command, error) {
	settings := struct{ apiURL, devUser, bearerToken string }{apiURL: defaultAPIURL}
	root := &cobra.Command{Use: "maintenance", Aliases: []string{"maint"}, Short: "Export, prune, and retain go-go-host operational state"}
	root.PersistentFlags().StringVar(&settings.apiURL, "api-url", defaultAPIURL, "go-go-host daemon API base URL")
	root.PersistentFlags().StringVar(&settings.devUser, "dev-user", "", "dev auth user subject header")
	root.PersistentFlags().StringVar(&settings.bearerToken, "bearer-token", "", "OIDC bearer token for non-dev auth mode")

	export := &cobra.Command{Use: "export", Short: "Export site metadata, SQLite DBs, and deployment bundles"}
	export.AddCommand(newExportMetadataCommand(&settings), newExportDBCommand(&settings), newExportBundleCommand(&settings))
	root.AddCommand(export, newPruneDeploymentsCommand(&settings), newAuditRetentionCommand(&settings))
	return root, nil
}

func newExportMetadataCommand(settings *struct{ apiURL, devUser, bearerToken string }) *cobra.Command {
	var siteID, output string
	cmd := &cobra.Command{Use: "metadata", Short: "Export org/site metadata as JSON", RunE: func(cmd *cobra.Command, args []string) error {
		if siteID == "" {
			return fmt.Errorf("--site-id is required")
		}
		if output == "" {
			output = fmt.Sprintf("%s-metadata.json", siteID)
		}
		return downloadToFile(settings.apiURL, fmt.Sprintf("/api/v1/sites/%s/export/metadata", siteID), settings.devUser, settings.bearerToken, output)
	}}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID")
	cmd.Flags().StringVarP(&output, "output-path", "o", "", "output path")
	return cmd
}

func newExportDBCommand(settings *struct{ apiURL, devUser, bearerToken string }) *cobra.Command {
	var siteID, output string
	cmd := &cobra.Command{Use: "db", Short: "Export a site's SQLite database", RunE: func(cmd *cobra.Command, args []string) error {
		if siteID == "" {
			return fmt.Errorf("--site-id is required")
		}
		if output == "" {
			output = fmt.Sprintf("%s.sqlite", siteID)
		}
		return downloadToFile(settings.apiURL, fmt.Sprintf("/api/v1/sites/%s/export/db", siteID), settings.devUser, settings.bearerToken, output)
	}}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID")
	cmd.Flags().StringVarP(&output, "output-path", "o", "", "output path")
	return cmd
}

func newExportBundleCommand(settings *struct{ apiURL, devUser, bearerToken string }) *cobra.Command {
	var deploymentID, output string
	cmd := &cobra.Command{Use: "bundle", Short: "Export an immutable deployment bundle", RunE: func(cmd *cobra.Command, args []string) error {
		if deploymentID == "" {
			return fmt.Errorf("--deployment-id is required")
		}
		if output == "" {
			output = fmt.Sprintf("%s-bundle.tar.gz", deploymentID)
		}
		return downloadToFile(settings.apiURL, fmt.Sprintf("/api/v1/deployments/%s/bundle", deploymentID), settings.devUser, settings.bearerToken, output)
	}}
	cmd.Flags().StringVar(&deploymentID, "deployment-id", "", "deployment ID")
	cmd.Flags().StringVarP(&output, "output-path", "o", "", "output path")
	return cmd
}

func newPruneDeploymentsCommand(settings *struct{ apiURL, devUser, bearerToken string }) *cobra.Command {
	var siteID, olderThan, statuses string
	var keepLatest int
	var dryRun bool
	cmd := &cobra.Command{Use: "prune-deployments", Short: "Prune old rejected/superseded deployment artifacts", RunE: func(cmd *cobra.Command, args []string) error {
		if siteID == "" {
			return fmt.Errorf("--site-id is required")
		}
		body := map[string]any{"olderThan": olderThan, "keepLatest": keepLatest, "dryRun": dryRun}
		if statuses != "" {
			body["statuses"] = strings.Split(statuses, ",")
		}
		var out any
		if err := postJSONWithAuth(settings.apiURL, fmt.Sprintf("/api/v1/sites/%s/deployments/prune", siteID), settings.devUser, settings.bearerToken, body, &out); err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
	}}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID")
	cmd.Flags().StringVar(&olderThan, "older-than", "0d", "RFC3339 timestamp or Nd duration")
	cmd.Flags().StringVar(&statuses, "statuses", "rejected,superseded", "comma-separated statuses")
	cmd.Flags().IntVar(&keepLatest, "keep-latest", 0, "keep latest N matching old deployments")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "preview candidates without deleting")
	return cmd
}

func newAuditRetentionCommand(settings *struct{ apiURL, devUser, bearerToken string }) *cobra.Command {
	var olderThan string
	var dryRun bool
	cmd := &cobra.Command{Use: "audit-retention", Short: "Delete old audit events as a platform admin", RunE: func(cmd *cobra.Command, args []string) error {
		if olderThan == "" {
			return fmt.Errorf("--older-than is required")
		}
		var out any
		if err := postJSONWithAuth(settings.apiURL, "/api/v1/admin/audit/retention", settings.devUser, settings.bearerToken, map[string]any{"olderThan": olderThan, "dryRun": dryRun}, &out); err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(out)
	}}
	cmd.Flags().StringVar(&olderThan, "older-than", "", "RFC3339 timestamp or Nd duration")
	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "preview without deleting")
	return cmd
}

func downloadToFile(apiURL, path, devUser, bearerToken, output string) error {
	base := strings.TrimRight(apiURL, "/")
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodGet, base+path, nil)
	if err != nil {
		return err
	}
	if devUser != "" {
		req.Header.Set("X-Go-Go-Host-User", devUser)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("GET %s%s: %w", base, path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GET %s%s: unexpected status %s: %s", base, path, resp.Status, strings.TrimSpace(string(data)))
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil && filepath.Dir(output) != "." {
		return err
	}
	return os.WriteFile(output, buf.Bytes(), 0o600)
}
