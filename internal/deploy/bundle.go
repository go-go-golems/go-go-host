package deploy

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const ManifestName = "go-go-host.json"

var SafeCapabilities = map[string]bool{"express": true, "ui.dsl": true, "database": true, "db": true, "time": true, "timer": true, "assets": true, "sqlite": true}

type Manifest struct {
	Name         string   `json:"name"`
	ScriptsDir   string   `json:"scriptsDir"`
	AssetsDir    string   `json:"assetsDir"`
	Entrypoint   string   `json:"entrypoint"`
	SmokePath    string   `json:"smokePath"`
	Capabilities []string `json:"capabilities"`
	AllowedPaths []string `json:"allowedPaths"`
	Channel      string   `json:"channel"`
}

type ValidationReport struct {
	Valid                 bool     `json:"valid"`
	Errors                []string `json:"errors,omitempty"`
	Warnings              []string `json:"warnings,omitempty"`
	Files                 int      `json:"files"`
	Bytes                 int64    `json:"bytes"`
	RequestedCapabilities []string `json:"requestedCapabilities,omitempty"`
	EffectiveCapabilities []string `json:"effectiveCapabilities,omitempty"`
}

func (r *ValidationReport) addError(format string, args ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, args...))
	r.Valid = false
}

type Options struct {
	MaxBytes        int64
	MaxFiles        int
	AllowedPaths    []string
	AllowedChannels []string
	Channel         string
	PolicyCaps      map[string]bool
	BundleFilename  string
}

type PreparedBundle struct {
	Manifest      Manifest
	ManifestJSON  []byte
	Report        ValidationReport
	ArchivePath   string
	UnpackedPath  string
	BundleSHA256  string
	EffectiveCaps []string
}

func ValidateAndStore(ctx context.Context, srcPath, archiveDest, unpackDest string, opts Options) (*PreparedBundle, error) {
	report := ValidationReport{Valid: true}
	if opts.MaxFiles <= 0 {
		opts.MaxFiles = 2000
	}
	policy := opts.PolicyCaps
	if policy == nil {
		policy = SafeCapabilities
	}
	if err := os.MkdirAll(filepath.Dir(archiveDest), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(unpackDest, 0o755); err != nil {
		return nil, err
	}
	archiveBytes, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}
	if opts.MaxBytes > 0 && int64(len(archiveBytes)) > opts.MaxBytes {
		report.addError("bundle archive exceeds quota: %d > %d bytes", len(archiveBytes), opts.MaxBytes)
	}
	sum := sha256.Sum256(archiveBytes)
	digest := hex.EncodeToString(sum[:])

	files, err := readArchive(srcPath)
	if err != nil {
		return &PreparedBundle{Report: report}, err
	}
	if len(files) > opts.MaxFiles {
		report.addError("bundle contains too many files: %d > %d", len(files), opts.MaxFiles)
	}
	var manifestData []byte
	for _, f := range files {
		if err := validatePath(f.Name); err != nil {
			report.addError("invalid path %q: %v", f.Name, err)
			continue
		}
		if !allowsPath(f.Name, opts.AllowedPaths) {
			report.addError("path %q is not allowed by deployment policy", f.Name)
		}
		report.Files++
		report.Bytes += int64(len(f.Data))
		if opts.MaxBytes > 0 && report.Bytes > opts.MaxBytes {
			report.addError("bundle unpacked bytes exceed quota: %d > %d", report.Bytes, opts.MaxBytes)
		}
		if f.Name == ManifestName {
			manifestData = f.Data
		}
	}
	var manifest Manifest
	if len(manifestData) == 0 {
		report.addError("missing %s manifest", ManifestName)
	} else if err := json.Unmarshal(manifestData, &manifest); err != nil {
		report.addError("invalid manifest JSON: %v", err)
	} else {
		validateManifest(&report, manifest, policy)
		if opts.Channel != "" && manifest.Channel != "" && manifest.Channel != opts.Channel {
			report.addError("manifest channel %q does not match requested channel %q", manifest.Channel, opts.Channel)
		}
		if len(opts.AllowedChannels) > 0 && manifest.Channel != "" && !containsString(opts.AllowedChannels, manifest.Channel) {
			report.addError("channel %q is not allowed by deployment policy", manifest.Channel)
		}
		pathPolicy := append([]string{}, opts.AllowedPaths...)
		pathPolicy = append(pathPolicy, manifest.AllowedPaths...)
		if len(pathPolicy) > 0 {
			for _, f := range files {
				if !allowsPath(f.Name, pathPolicy) {
					report.addError("path %q is not allowed by manifest/deployment policy", f.Name)
				}
			}
		}
	}
	effective := intersectCapabilities(manifest.Capabilities, policy)
	report.RequestedCapabilities = manifest.Capabilities
	report.EffectiveCapabilities = effective
	manifestJSON, _ := json.Marshal(manifest)
	if !report.Valid {
		return &PreparedBundle{Manifest: manifest, ManifestJSON: manifestJSON, Report: report, BundleSHA256: digest, EffectiveCaps: effective}, nil
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err := os.WriteFile(archiveDest, archiveBytes, 0o644); err != nil {
		return nil, err
	}
	if err := unpackFiles(files, unpackDest); err != nil {
		return nil, err
	}
	return &PreparedBundle{Manifest: manifest, ManifestJSON: manifestJSON, Report: report, ArchivePath: archiveDest, UnpackedPath: unpackDest, BundleSHA256: digest, EffectiveCaps: effective}, nil
}

func validateManifest(report *ValidationReport, manifest Manifest, policy map[string]bool) {
	if strings.TrimSpace(manifest.ScriptsDir) == "" {
		report.addError("manifest scriptsDir is required")
	}
	for _, path := range []string{manifest.ScriptsDir, manifest.AssetsDir, manifest.Entrypoint} {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if err := validatePath(path); err != nil {
			report.addError("invalid manifest path %q: %v", path, err)
		}
	}
	for _, cap := range manifest.Capabilities {
		if !policy[cap] {
			report.addError("capability %q is not permitted by site policy", cap)
		}
	}
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func intersectCapabilities(requested []string, policy map[string]bool) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, cap := range requested {
		if policy[cap] && !seen[cap] {
			out = append(out, cap)
			seen[cap] = true
		}
	}
	return out
}

type archiveFile struct {
	Name string
	Data []byte
}

func readArchive(path string) ([]archiveFile, error) {
	if strings.HasSuffix(strings.ToLower(path), ".zip") {
		return readZip(path)
	}
	return readTarGz(path)
}

func readTarGz(path string) ([]archiveFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	out := []archiveFile{}
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag == tar.TypeSymlink || h.Typeflag == tar.TypeLink {
			return nil, fmt.Errorf("unsafe link entry %q", h.Name)
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		out = append(out, archiveFile{Name: normalizeArchiveName(h.Name), Data: data})
	}
	return out, nil
}

func readZip(path string) ([]archiveFile, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	out := []archiveFile{}
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if f.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("unsafe symlink entry %q", f.Name)
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(rc)
		_ = rc.Close()
		if readErr != nil {
			return nil, readErr
		}
		out = append(out, archiveFile{Name: normalizeArchiveName(f.Name), Data: data})
	}
	return out, nil
}

func normalizeArchiveName(path string) string {
	return filepath.ToSlash(filepath.Clean(filepath.ToSlash(path)))
}

func validatePath(path string) error {
	if path == "" {
		return errors.New("empty path")
	}
	path = filepath.ToSlash(path)
	if strings.HasPrefix(path, "/") || filepath.IsAbs(path) {
		return errors.New("absolute paths are forbidden")
	}
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "." || strings.HasPrefix(clean, "../") || clean == ".." || strings.Contains(clean, "/../") {
		return errors.New("parent traversal is forbidden")
	}
	for _, part := range strings.Split(clean, "/") {
		if part == "" || part == "." || part == ".." {
			return errors.New("unsafe path component")
		}
		if strings.HasPrefix(part, ".") && part != ".well-known" {
			return fmt.Errorf("hidden metadata component %q is forbidden", part)
		}
	}
	return nil
}

func allowsPath(path string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, pattern := range allowed {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		if pattern == "**" || pattern == path {
			return true
		}
		if strings.HasSuffix(pattern, "/**") && strings.HasPrefix(path, strings.TrimSuffix(pattern, "**")) {
			return true
		}
		if ok, _ := filepath.Match(pattern, path); ok {
			return true
		}
	}
	return false
}

func unpackFiles(files []archiveFile, dest string) error {
	if err := os.RemoveAll(dest); err != nil {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	for _, f := range files {
		if err := validatePath(f.Name); err != nil {
			return err
		}
		target := filepath.Join(dest, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("unsafe unpack target %q", f.Name)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, f.Data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
