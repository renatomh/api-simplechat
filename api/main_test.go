package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/renatomh/api-simplechat/db/sqlc"
	"github.com/renatomh/api-simplechat/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

// TestMain is the managing point for all tests to be run in the package
func TestMain(m *testing.M) {
	// Here, we set gin to test mode, in order to avoid dbug output overload
	gin.SetMode(gin.TestMode)

	// Running the tests
	os.Exit(m.Run())
}
