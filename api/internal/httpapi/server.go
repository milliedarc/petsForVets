// Package httpapi turns HTTP requests into database queries and JSON responses.
package httpapi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"petsforvets/api/internal/db"
)

// Server holds what every handler needs. Handlers are methods on it, so they
// can reach the database through s.queries.
type Server struct {
	pool         *pgxpool.Pool // needed to start transactions
	queries      *db.Queries
	secureCookie bool // true in production (HTTPS only), false on localhost
}

func NewServer(pool *pgxpool.Pool, secureCookie bool) *Server {
	return &Server{pool: pool, queries: db.New(pool), secureCookie: secureCookie}
}

// inTx runs fn inside a database transaction: if fn returns an error, every
// change it made is rolled back; otherwise they're all committed together.
// The PDO equivalent is beginTransaction() / commit() / rollBack().
func (s *Server) inTx(ctx context.Context, fn func(q *db.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	// Rollback after a successful Commit is a harmless no-op, so deferring it
	// covers every early return and panic.
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(s.queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
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
	mux.HandleFunc("PUT /api/auth/me/app-mode", s.requireUser(s.handleUpdateAppMode))
	mux.HandleFunc("POST /api/auth/change-password", s.requireUser(s.handleChangePassword))

	mux.HandleFunc("GET /api/pets", s.requireUser(s.handleListPets))
	mux.HandleFunc("POST /api/pets", s.requireUser(s.handleCreatePet))
	mux.HandleFunc("GET /api/pets/{id}", s.requireUser(s.handleGetPet))
	mux.HandleFunc("PUT /api/pets/{id}", s.requireUser(s.handleUpdatePet))
	mux.HandleFunc("DELETE /api/pets/{id}", s.requireUser(s.handleDeletePet))

	mux.HandleFunc("GET /api/pets/{id}/prescriptions", s.requireUser(s.handleListPrescriptions))
	mux.HandleFunc("GET /api/pets/{id}/prescription-requests", s.requireUser(s.handleListPrescriptionRequests))
	mux.HandleFunc("POST /api/prescriptions/{id}/requests", s.requireUser(s.handleCreatePrescriptionRequest))
	mux.HandleFunc("DELETE /api/prescription-requests/{id}", s.requireUser(s.handleCancelPrescriptionRequest))

	// Clinic staff only (requireStaff also checks users.veterinary).
	mux.HandleFunc("GET /api/clinic/summary", s.requireStaff(s.handleClinicSummary))
	mux.HandleFunc("GET /api/clinic/prescription-requests", s.requireStaff(s.handleListClinicRequests))
	mux.HandleFunc("PATCH /api/clinic/prescription-requests/{id}", s.requireStaff(s.handleUpdateClinicRequest))
	mux.HandleFunc("GET /api/clinic/pets", s.requireStaff(s.handleListClinicPets))
	mux.HandleFunc("GET /api/clinic/pets/{id}", s.requireStaff(s.handleGetClinicPet))
	mux.HandleFunc("GET /api/clinic/pets/{id}/prescriptions", s.requireStaff(s.handleListClinicPetPrescriptions))
	mux.HandleFunc("POST /api/clinic/pets/{id}/prescriptions", s.requireStaff(s.handleCreatePrescription))
	mux.HandleFunc("GET /api/clinic/medicines", s.requireStaff(s.handleListMedicines))

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
