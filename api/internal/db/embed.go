package db

import "embed"

// Migrations holds the .sql migration files, compiled into the binary so the
// server can bring the database schema up to date when it starts.
//
//go:embed migrations/*.sql
var Migrations embed.FS
