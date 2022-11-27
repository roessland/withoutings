package migration

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"os"
)

func Run(db *sql.DB) {
	// Load migration files embedded in binary
	migrationData, err := iofs.New(FS, ".")
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
}
