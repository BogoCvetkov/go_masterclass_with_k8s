package main

import (
	"database/sql"
	"os"

	"github.com/BogoCvetkov/go_mastercalss/api"
	validator "github.com/BogoCvetkov/go_mastercalss/api/controller/validators"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	g "github.com/BogoCvetkov/go_mastercalss/grpc"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	// Load Config
	config := config.LoadConfig()

	// Config logger
	if config.ENV == "DEV" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Initialize DB connection and Store
	conn, err := sql.Open(config.DB_DRIVER, config.DB_URL)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed connecting to DB")
	}

	store := db.NewStore(conn)

	// Initialize new server
	server := api.NewServer(store, config)

	defer server.GetAsync().Close()

	// Attach middlewares used everywhere
	server.AttachGlobalMiddlewares()

	// Attach route
	server.AttachRoutes()

	// Init custom validators
	validator.RegisterValidation()

	// Init GRPC server in parallel
	gserver := g.NewServer(store, config)
	go gserver.Start(config.GRPC_PORT)

	// GRPC gateway
	go gserver.StartHttpGateway(config.GRPC_GATEWAY_PORT)

	// Start listening
	server.Start(config.PORT)
}
