package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/internal/logging"
	"github.com/roessland/withoutings/internal/migrations"
	"strconv"
	"time"
)

// TestDatabase contains an empty database with all migrations applied.
type TestDatabase struct {
	*pgxpool.Pool // wotrw
	postgresPool  *pgxpool.Pool
	dbName        string
}

// New creates a new test database.
func New(ctx context.Context) TestDatabase {
	logger := logging.MustGetLoggerFromContext(ctx)

	// Connect to postgres using socket/trust with default user (probably <username> or postgres)
	logger.Debugf("Connecting to template1")
	postgresPool, err := pgxpool.Connect(ctx, "postgres://?host=/tmp&database=template1")
	if err != nil {
		panic(err)
	}

	// Create new test database as postgres user
	dbName := "wot_test_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	logger.Debugf("Creating database " + dbName)
	_, err = postgresPool.Exec(ctx, `
	create database `+dbName+`
	owner wotsa
	template template0
	encoding 'utf8'
	lc_collate = 'C';`)
	if err != nil {
		panic(err)
	}

	// Connect to test DB as superuser
	logger.Debugf("Connecting as wotsa")
	wotsaConn, err := sql.Open("pgx", "postgres://?host=/tmp&search_path=wot&user=wotsa&database="+dbName)
	if err != nil {
		panic(err)
	}

	// Make sure wot schema exists, or golang-migrate will fail with a null current_schema error.
	logger.Debugf("Creating schema")
	_, err = wotsaConn.Exec(`create schema if not exists wot;`)
	if err != nil {
		panic(err)
	}

	// Run migrations on test DB as superuser
	logger.Debugf("Running migrations")
	migrations.Run(wotsaConn)

	// Connect to test DB using wotrw user
	wotrwPool, err := pgxpool.Connect(ctx, "postgres://?host=/tmp&&user=wotrw&database="+dbName)
	if err != nil {
		panic(err)
	}

	// Return wotrw user connection
	return TestDatabase{
		Pool:         wotrwPool,
		postgresPool: postgresPool,
		dbName:       dbName,
	}
}

// Drop drops the test database.
func (tdb *TestDatabase) Drop(ctx context.Context) {
	logger := logging.MustGetLoggerFromContext(ctx)

	// wotsa/Superadmin connection was already closed, but the idle connection can still be
	// around, so we need to drop with force.

	// Close readwrite connection.
	logger.Debugf("Closing RW pool")
	tdb.Pool.Close()

	// Drop database using postgres connection.
	logger.Debugf("Dropping DB")
	_, err := tdb.postgresPool.Exec(context.Background(), fmt.Sprintf("drop database %s with (force);", tdb.dbName))
	if err != nil {
		panic(err)
	}

	// Close postgres connection
	logger.Debugf("Closing postgres connection")
	tdb.postgresPool.Close()
}
