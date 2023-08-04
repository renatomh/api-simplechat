package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestMain is the managing point for all tests to be run in the package
func TestMain(m *testing.M) {
	// Here, we set gin to test mode, in order to avoid dbug output overload
	gin.SetMode(gin.TestMode)

	// Running the tests
	os.Exit(m.Run())
}
