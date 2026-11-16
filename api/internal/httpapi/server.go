// Package httpapi turns HTTP requests into database queries and JSON responses.
package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"petsforvets/api/internal/db"
)

// Server holds what every handler needs. Handlers are methods on it, so they
// can reach the database through s.queries.
type Server struct {
	queries      *db.Queries
	secureCookie bool // true in production (HTTPS only), false on localhost
}

func NewServer(queries *db.Queries, secureCookie bool) *Server {
	return &Server{queries: queries, secureCookie: secureCookie}
}

// Routes is the routing table: "METHOD /path" -> handler function.
// Since Go 1.22 the standard library router understands methods and
// {placeholders}, so no framework is needed.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", s.handleRegister)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("GET /api/auth/me", s.requireUser(s.handleMe))

	mux.HandleFunc("GET /api/pets", s.requireUser(s.handleListPets))
	mux.HandleFunc("POST /api/pets", s.requireUser(s.handleCreatePet))
	mux.HandleFunc("GET /api/pets/{id}", s.requireUser(s.handleGetPet))
	mux.HandleFunc("PUT /api/pets/{id}", s.requireUser(s.handleUpdatePet))
	mux.HandleFunc("DELETE /api/pets/{id}", s.requireUser(s.handleDeletePet))

	mux.HandleFunc("GET /api/species", s.requireUser(s.handleListSpecies))
	mux.HandleFunc("GET /api/breeds", s.requireUser(s.handleListBreeds))

	// Middleware wraps the router, outermost first:
	//  - CrossOriginProtection rejects POST/PUT/DELETE requests coming from
	//    other websites (CSRF protection), using headers browsers send.
	//  - logRequests prints one line per request.
	return http.NewCrossOriginProtection().Handler(logRequests(mux))
}

// statusRecorder remembers the status code a handler wrote, for logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start).Round(time.Microsecond),
		)
	})
}
