package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/roessland/withoutings/pkg/logging"
	"github.com/roessland/withoutings/pkg/migration"
	"strconv"
	"time"
)

// TestDatabase contains an empty database with all migrations applied.
type TestDatabase struct {
	*pgxpool.Pool               // Connection with wotrw user (read/write)
	postgresPool  *pgxpool.Pool // Connection with postgres user (superuser)
	dbName        string
}

// New creates a new test database.
func New(ctx context.Context) TestDatabase {
	logger := logging.MustGetLoggerFromContext(ctx)

	// Connect to postgres using socket/trust with current user
	// (<myuser> on localhost, "runner" in CI)
	logger.Debugf("Connecting to template1")
	postgresPool, err := pgxpool.Connect(ctx, "postgres://?host=localhost&password=postgres&database=template1")
	if err != nil {
		panic(err)
	}

	// Create role for superadmin user
	logger.Debugf("Creating wotsa user ")
	_, _ = postgresPool.Exec(ctx, `
		create role wotsa
		password 'wotsa'
		login;`) // Ignore error -- user might already exist

	// Create role for readwrite user
	logger.Debugf("Creating wotrw user")
	_, _ = postgresPool.Exec(ctx, `
		create role wotrw
		password 'wotrw'
		login;`) // Ignore error -- user might already exist

	// Create new test database as postgres user
	dbName := "wot_test_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	logger.Infof("Creating database " + dbName)
	_, err = postgresPool.Exec(ctx, `
		create database `+dbName+`
		owner wotsa
		template template0
		encoding 'utf8'
		lc_collate = 'C';`)
	if err != nil {
		panic(err)
	}

	// Set search path for all _future_ connections
	logger.Debugf("Setting search path")
	_, err = postgresPool.Exec(ctx, fmt.Sprintf(`
		alter database %s set search_path to wot;`, dbName))
	if err != nil {
		panic(err)
	}

	// Grant temp table permissions
	logger.Debugf("Granting temp table permission to wotrw")
	_, err = postgresPool.Exec(ctx, fmt.Sprintf(`
		grant temporary on database %s to wotrw;`, dbName))
	if err != nil {
		panic(err)
	}

	// Connect to test DB as superuser
	logger.Debugf("Connecting as wotsa")
	wotsaConn, err := sql.Open("pgx", "postgres://?host=localhost&search_path=wot&user=wotsa&password=wotsa&database="+dbName)
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
	migration.Run(wotsaConn)

	// Connect to test DB using wotrw user
	wotrwPool, err := pgxpool.Connect(ctx, "postgres://?host=localhost&user=wotrw&password=wotrw&database="+dbName)
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
