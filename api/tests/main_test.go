package api

import (
	"os"
	"testing"
	"time"

	"github.com/BogoCvetkov/go_mastercalss/api"
	validator "github.com/BogoCvetkov/go_mastercalss/api/controller/validators"
	"github.com/BogoCvetkov/go_mastercalss/config"
	db "github.com/BogoCvetkov/go_mastercalss/db"
	testutil "github.com/BogoCvetkov/go_mastercalss/db/tests"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T, store db.IStore) *api.Server {
	config := &config.Config{
		TOKEN_SECRET:   testutil.RandomString(32),
		TOKEN_DURATION: time.Minute,
	}

	server := api.NewServer(store, config)

	// Attach middlewares used everywhere
	server.AttachGlobalMiddlewares()

	// Attach route
	server.AttachRoutes()

	// Init custom validators
	validator.RegisterValidation()

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
