package docfs

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"

	hostdoc "github.com/go-go-golems/go-go-host/cmd/go-go-host/doc"
	agentdoc "github.com/go-go-golems/go-go-host/cmd/go-go-host-agent/doc"
)

// DocEntry is the JSON shape returned by the docs API.
type DocEntry struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Short   string `json:"short"`
	Section string `json:"section"`
	Source  string `json:"source"` // "host" or "agent"
	Body    string `json:"body,omitempty"`
}

var catalogue []DocEntry
var seenSlugs = map[string]struct{}{}

func init() {
	type source struct {
		fs     http.FileSystem
		name   string
	}
	sources := []source{
		{http.FS(hostdoc.DocFS()), "host"},
		{http.FS(agentdoc.DocFS()), "agent"},
	}

	for _, src := range sources {
		f, err := src.fs.Open("/")
		if err != nil {
			continue
		}
		entries, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			rawBytes, err := readAll(src.fs, e.Name())
			if err != nil {
				continue
			}
			rawStr := string(rawBytes)
			fm := parseFrontmatter(rawStr)
			body := stripFrontmatter(rawStr)
			slug := fm["Slug"]
			if slug == "" {
				slug = strings.TrimSuffix(e.Name(), ".md")
			}
			// Track slug per source to detect collisions.
			// If two sources produce the same slug, the second one
			// gets a source suffix (e.g. "agent-guide" -> "agent-guide-agent").
			slugKey := slug
			if _, exists := seenSlugs[slugKey]; exists {
				slug = slug + "-" + src.name
			}
			seenSlugs[slugKey] = struct{}{}
			catalogue = append(catalogue, DocEntry{
				Slug:    slug,
				Title:   fm["Title"],
				Short:   fm["Short"],
				Section: fm["SectionType"],
				Source:  src.name,
				Body:    body,
			})
		}
	}

	sort.Slice(catalogue, func(i, j int) bool {
		order := map[string]int{"Tutorial": 0, "GeneralTopic": 1}
		oI, oJ := order[catalogue[i].Section], order[catalogue[j].Section]
		if oI != oJ {
			return oI < oJ
		}
		return catalogue[i].Title < catalogue[j].Title
	})
}

func readAll(fs http.FileSystem, name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// HandleListDocs returns the doc catalogue (without body content).
func HandleListDocs(w http.ResponseWriter, _ *http.Request) {
	entries := make([]DocEntry, len(catalogue))
	for i, d := range catalogue {
		entries[i] = DocEntry{
			Slug:    d.Slug,
			Title:   d.Title,
			Short:   d.Short,
			Section: d.Section,
			Source:  d.Source,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

// HandleGetDoc returns a single doc by slug (with body content).
func HandleGetDoc(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	for _, d := range catalogue {
		if d.Slug == slug {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(d)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "doc not found"})
}

func parseFrontmatter(raw string) map[string]string {
	fm := map[string]string{}
	if !strings.HasPrefix(raw, "---\n") {
		return fm
	}
	end := strings.Index(raw[4:], "\n---")
	if end < 0 {
		return fm
	}
	for _, line := range strings.Split(raw[4:4+end], "\n") {
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"'`)
		fm[key] = val
	}
	return fm
}

func stripFrontmatter(raw string) string {
	if !strings.HasPrefix(raw, "---\n") {
		return raw
	}
	end := strings.Index(raw[4:], "\n---")
	if end < 0 {
		return raw
	}
	return raw[4+end+4:]
}
