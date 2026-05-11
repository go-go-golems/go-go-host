package httpapi

import "net/http"

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
	if rec.statusCode == http.StatusNotFound && !rec.wroteBody {
		h.fallback.ServeHTTP(w, r)
	}
}

type fallbackRecorder struct {
	http.ResponseWriter
	statusCode int
	wroteBody  bool
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
		return len(b), nil
	}
	r.wroteBody = true
	return r.ResponseWriter.Write(b)
}
