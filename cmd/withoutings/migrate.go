package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/roessland/withoutings/pkg/migration"
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

	// Run migrations
	migration.Run(db)
	fmt.Fprintf(os.Stdout, "Migrations complete\n")
}
