// The API server. Run from the api/ folder with:  go run ./cmd/server
package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"petsforvets/api/internal/db"
	"petsforvets/api/internal/httpapi"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

// run does the real work. It returns errors instead of exiting so that the
// deferred cleanup (closing the database pool) always happens.
func run() error {
	// ctx is cancelled when you press Ctrl+C, which starts a clean shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseURL := getenv("DATABASE_URL",
		"postgres://petsforvets:petsforvets@localhost:5432/petsforvets?sslmode=disable")
	addr := getenv("ADDR", "localhost:8080")
	secureCookie := os.Getenv("COOKIE_SECURE") == "true"

	// A pool keeps a few database connections open and shares them between
	// requests (PHP opens a new PDO connection on every request instead).
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("creating database pool: %w", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("connecting to database (is `docker-compose up -d` running?): %w", err)
	}

	if err := migrate(ctx, pool); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	server := httpapi.NewServer(db.New(pool), secureCookie)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Each request runs in its own goroutine automatically. This extra one just
	// waits for Ctrl+C and then lets in-flight requests finish.
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("shutting down", "err", err)
		}
	}()

	slog.Info("listening", "url", "http://"+addr)
	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-shutdownDone
	slog.Info("shut down cleanly")
	return nil
}

// migrate applies any migrations in internal/db/migrations that haven't run yet.
func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrations, err := fs.Sub(db.Migrations, "migrations")
	if err != nil {
		return err
	}

	sqlDB := stdlib.OpenDBFromPool(pool) // goose speaks database/sql, not pgx
	defer sqlDB.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations)
	if err != nil {
		return err
	}
	results, err := provider.Up(ctx)
	for _, result := range results {
		slog.Info("applied migration", "file", result.Source.Path)
	}
	return err
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
