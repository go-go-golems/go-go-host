// build-web is the reproducible Dagger/local build pipeline for the React
// dashboard. It builds web/admin with pnpm and copies dist/ into
// internal/webadmin/dist for go:embed.
//
// Usage:
//
//	go run ./cmd/build-web
//	BUILD_WEB_LOCAL=1 go run ./cmd/build-web
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

const (
	defaultBuilderImage = "node:22"
	defaultPNPMVersion  = "10.13.1"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	if os.Getenv("BUILD_WEB_LOCAL") == "1" {
		return runLocal(repoRoot)
	}
	if err := runDagger(ctx, repoRoot); err != nil {
		if errors.Is(err, errDaggerUnavailable) {
			fmt.Fprintln(os.Stderr, "dagger unavailable (no Docker/Dagger?), falling back to local pnpm")
			return runLocal(repoRoot)
		}
		return err
	}
	return nil
}

var errDaggerUnavailable = errors.New("dagger: engine not reachable")

func runDagger(ctx context.Context, repoRoot string) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("%w: %v", errDaggerUnavailable, err)
	}
	defer func() { _ = client.Close() }()

	webDir := filepath.Join(repoRoot, "web", "admin")
	pnpmVersion := getenv("WEB_PNPM_VERSION", readPNPMVersion(webDir))
	if pnpmVersion == "" {
		pnpmVersion = defaultPNPMVersion
	}
	builderImage := getenv("WEB_BUILDER_IMAGE", defaultBuilderImage)

	source := client.Host().Directory(webDir, dagger.HostDirectoryOpts{Exclude: []string{"dist", "storybook-static", "node_modules", ".git"}})
	pnpmStore := client.CacheVolume("go-go-host-admin-pnpm-store")
	pathEnv := "/pnpm:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

	container := client.Container().
		From(builderImage).
		WithEnvVariable("PNPM_HOME", "/pnpm").
		WithEnvVariable("PATH", pathEnv).
		WithMountedCache("/pnpm/store", pnpmStore).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-lc", "corepack enable && corepack prepare pnpm@" + pnpmVersion + " --activate"}).
		WithExec([]string{"pnpm", "install", "--frozen-lockfile", "--prefer-offline"}).
		WithExec([]string{"pnpm", "run", "build"})

	tmpDir, err := os.MkdirTemp("", "go-go-host-admin-dist-")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	if _, err := container.Directory("/src/dist").Export(ctx, tmpDir); err != nil {
		return fmt.Errorf("export dist: %w", err)
	}
	if err := copyDistToEmbed(repoRoot, tmpDir); err != nil {
		return err
	}
	log.Printf("Successfully exported web/admin dist to internal/webadmin/dist (via Dagger)")
	return nil
}

func runLocal(repoRoot string) error {
	webDir := filepath.Join(repoRoot, "web", "admin")
	if err := runPNPM(webDir, "install", "--frozen-lockfile", "--prefer-offline"); err != nil {
		return fmt.Errorf("pnpm install (local): %w", err)
	}
	if err := runPNPM(webDir, "run", "build"); err != nil {
		return fmt.Errorf("pnpm run build (local): %w", err)
	}
	if err := copyDistToEmbed(repoRoot, filepath.Join(webDir, "dist")); err != nil {
		return err
	}
	log.Printf("Successfully exported web/admin dist to internal/webadmin/dist (local pnpm)")
	return nil
}

func copyDistToEmbed(repoRoot, src string) error {
	dst := filepath.Join(repoRoot, "internal", "webadmin", "dist")
	if err := recreate(dst); err != nil {
		return fmt.Errorf("recreate embed dist: %w", err)
	}
	if err := copyTree(src, dst); err != nil {
		return fmt.Errorf("copy to embed dist: %w", err)
	}
	return nil
}

func readPNPMVersion(webDir string) string {
	data, err := os.ReadFile(filepath.Join(webDir, "package.json"))
	if err != nil {
		return ""
	}
	var pkg struct {
		PackageManager string `json:"packageManager"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}
	return strings.TrimPrefix(pkg.PackageManager, "pnpm@")
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func runPNPM(dir string, args ...string) error {
	// #nosec G204 -- command is fixed to pnpm; only controlled subcommands/arguments are supplied.
	cmd := exec.Command("pnpm", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func recreate(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, 0o755)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		// #nosec G122 -- filepath.WalkDir provides p from the trusted build output tree.
		in, err := os.Open(p)
		if err != nil {
			return err
		}
		defer func() { _ = in.Close() }()
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, in)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
}
