package main

import (
	"database/sql"
	"log"

	"github.com/BogoCvetkov/go_mastercalss/api"
	validator "github.com/BogoCvetkov/go_mastercalss/controller/validators"
	"github.com/BogoCvetkov/go_mastercalss/db"
	_ "github.com/lib/pq"
)

func main() {

	// Load Config
	config := LoadConfig()

	// Initialize DB connection and Store
	conn, err := sql.Open(config.DBDriver, config.DB)

	if err != nil {
		log.Fatal("Failed connecting to DB", err)
	}

	store := db.NewStore(conn)

	// Initialize new server
	server := api.NewServer(store)

	// Attach middlewares used everywhere
	server.AttachGlobalMiddlewares()

	// Attach route
	server.AttachRoutes()

	// Init custom validators
	validator.RegisterValidation()

	// Start listening
	server.Start(config.Port)
}
