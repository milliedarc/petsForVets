package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// writeJSON sends v as a JSON response with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writing JSON response", "err", err)
	}
}

// writeError sends {"error": message}. The frontend shows or logs this.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// serverError is for unexpected failures (database down, bugs). The real error
// goes to the server log; the client only gets a generic message so we don't
// leak internals.
func serverError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal error", "method", r.Method, "path", r.URL.Path, "err", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

// readJSON decodes the request body into dst. Bodies over 1 MB and unknown
// fields are rejected, so typos in the frontend show up as errors.
func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}
