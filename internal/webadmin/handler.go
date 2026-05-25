package webadmin

import (
	"bytes"
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

//go:embed all:dist
var dashboardFS embed.FS

// NewHandler serves the embedded Vite dashboard under an already-stripped
// prefix. Nested SPA routes fall back to index.html while asset paths are served
// directly from the embedded dist directory.
func NewHandler() http.Handler {
	dist, err := fs.Sub(dashboardFS, "dist")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		})
	}
	index, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		})
	}
	files := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if requestPath == "" || requestPath == "/" {
			serveIndex(w, r, index)
			return
		}
		clean := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))
		if strings.Contains(path.Base(clean), ".") {
			r.URL.Path = clean
			files.ServeHTTP(w, r)
			return
		}
		serveIndex(w, r, index)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, index []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(index))
}

// NewPlaceholderHandler is kept as a compatibility alias for older call sites.
func NewPlaceholderHandler() http.Handler { return NewHandler() }
