package httpapi

import (
	"net/http"
	"strconv"
)

func parseInt32QueryParam(r *http.Request, name string) int {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return 0
	}
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0
	}
	return int(value)
}
