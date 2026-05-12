package httpapi

import (
	"bytes"
	"net/http"
	"strings"
)

type fallbackHandler struct {
	primary  http.Handler
	fallback http.Handler
}

func withFallback(primary, fallback http.Handler) http.Handler {
	if fallback == nil {
		return primary
	}
	return fallbackHandler{primary: primary, fallback: fallback}
}

func (h fallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := &fallbackRecorder{ResponseWriter: w, statusCode: http.StatusOK}
	h.primary.ServeHTTP(rec, r)
	if rec.statusCode == http.StatusNotFound && shouldFallbackToHostedSite(r.URL.Path) {
		h.fallback.ServeHTTP(w, r)
		return
	}
	if rec.statusCode == http.StatusNotFound {
		rec.flushNotFound()
	}
}

func shouldFallbackToHostedSite(path string) bool {
	for _, prefix := range []string{"/api/", "/app/", "/admin/"} {
		if strings.HasPrefix(path, prefix) {
			return false
		}
	}
	switch path {
	case "/api", "/app", "/admin", "/healthz", "/readyz":
		return false
	default:
		return true
	}
}

type fallbackRecorder struct {
	http.ResponseWriter
	statusCode int
	wroteBody  bool
	notFound   bytes.Buffer
}

func (r *fallbackRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	if statusCode != http.StatusNotFound {
		r.ResponseWriter.WriteHeader(statusCode)
	}
}

func (r *fallbackRecorder) Write(b []byte) (int, error) {
	if r.statusCode == http.StatusNotFound {
		r.wroteBody = true
		return r.notFound.Write(b)
	}
	r.wroteBody = true
	return r.ResponseWriter.Write(b)
}

func (r *fallbackRecorder) flushNotFound() {
	r.ResponseWriter.WriteHeader(http.StatusNotFound)
	if r.notFound.Len() > 0 {
		_, _ = r.ResponseWriter.Write(r.notFound.Bytes())
	}
}
