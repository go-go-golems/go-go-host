package webadmin

import (
	"fmt"
	"net/http"
)

// NewPlaceholderHandler returns a minimal dashboard placeholder. The real
// React/Vite dashboard will be embedded here in the dashboard phases.
func NewPlaceholderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<!doctype html>
<html>
  <head><title>go-go-host</title></head>
  <body>
    <h1>go-go-host</h1>
    <p>Dashboard placeholder. User dashboard and platform admin console will be embedded here.</p>
  </body>
</html>`)
	})
}
