package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/renatomh/api-simplechat/util"
)

var testQueries *Queries
var testDB *sql.DB

// TestMain is the managing point for all tests to be run in the package
func TestMain(m *testing.M) {
	var err error

	// Loading config file from root dir
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Connecting to the test database
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	// If an error is returned
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	// Initializing the test queries with the established connection
	testQueries = New(testDB)

	// Running the tests
	os.Exit(m.Run())
}
