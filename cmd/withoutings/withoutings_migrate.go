package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/roessland/withoutings/migrations"
	"os"
)

func withoutingsMigrate() {
	// Get connection string
	connectionString := os.Getenv("WOT_DATABASE_URL_SA")
	if connectionString == "" {
		fmt.Fprintf(os.Stderr, "WOT_DATABASE_URL_SA environment variable was empty\n")
		os.Exit(1)
	}

	// Connect to DB as superuser
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Load migration files embedded in binary
	migrationData, err := iofs.New(migrations.FS, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable load migration data: %v\n", err)
		os.Exit(1)
	}

	// Create migration migrationDriver
	migrationDriver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create migration migrationDriver: %v\n", err)
		os.Exit(1)
	}

	// Create migration instance
	migrationInstance, err := migrate.NewWithInstance(
		"iofs", migrationData,
		"postgres", migrationDriver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create migration instance: %v\n", err)
		os.Exit(1)
	}

	// Migrate
	err = migrationInstance.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Fprintf(os.Stderr, "Unable to run migrations:: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Migrations complete\n")
}
