package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	connString = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// We need a Queries object to test the methods implemented in the sqlc generated code.
var testQueries *Queries

// TestMain is the entry point for the test suite. It sets up the test environment
// and runs the tests.
// It is called by the Go testing framework before any tests are run.
// It is typically used to set up any necessary resources, such as database connections,
// configuration, or test data.
func TestMain(m *testing.M) {
	ctx := context.Background()

	// We can even use pgxpool.Pool to create a connection pool.
	// It automatically manages the connections for us, it's thread-safe and ideal for concurrent requests.
	// However, for simplicity, we are using pgx.Connect here.
	// conn, err := pgxpool.Connect(ctx, connString)
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatal("error while connecting to the db:", err)
	}
	defer conn.Close(ctx)

	testQueries = New(conn)

	// m.Run() is a function that runs the tests in the current package.
	// It returns an integer exit code, which is typically 0 for success and 1 for failure.
	// The exit code is passed to os.Exit, which terminates the program with that exit code.
	// This is useful for running tests in a continuous integration (CI) environment,
	// where the exit code can be used to determine if the tests passed or failed.
	os.Exit(m.Run())
}
