package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_chat?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain is the managing point for all tests to be run in the package
func TestMain(m *testing.M) {
	var err error

	// Connecting to the test database
	testDB, err = sql.Open(dbDriver, dbSource)
	// If an error is returned
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	// Initializing the test queries with the established connection
	testQueries = New(testDB)

	// Running the tests
	os.Exit(m.Run())
}
