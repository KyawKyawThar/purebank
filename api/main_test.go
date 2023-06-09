package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	db "purebank/db/sqlc"
	"purebank/db/util"
	"testing"
	"time"
)

func newTestServer(t *testing.T, s db.Store) *Server {

	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration:     time.Minute,
	}

	server, err := NewServer(config, s, nil)

	require.NoError(t, err)
	return server

}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
