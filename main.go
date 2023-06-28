package main

import (
	"database/sql"
	"log"

	"github.com/BogoCvetkov/go_mastercalss/api"
	validator "github.com/BogoCvetkov/go_mastercalss/api/controller/validators"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	g "github.com/BogoCvetkov/go_mastercalss/grpc"
	_ "github.com/lib/pq"
)

func main() {

	// Load Config
	config := config.LoadConfig()

	// Initialize DB connection and Store
	conn, err := sql.Open(config.DBDriver, config.DB)

	if err != nil {
		log.Fatal("Failed connecting to DB", err)
	}

	store := db.NewStore(conn)

	// Initialize new server
	server := api.NewServer(store, config)

	// Attach middlewares used everywhere
	server.AttachGlobalMiddlewares()

	// Attach route
	server.AttachRoutes()

	// Init custom validators
	validator.RegisterValidation()

	// Init GRP server in parallel
	gserver := g.NewServer(store, config)
	go gserver.Start(config.GRPCPort)

	// GRPC gateway
	go gserver.StartHttpGateway(config.GRPCGatewayPort)

	// Start listening
	server.Start(config.Port)
}
