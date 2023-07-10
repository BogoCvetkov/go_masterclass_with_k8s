package main

import (
	"database/sql"
	"fmt"

	"github.com/BogoCvetkov/go_mastercalss/async"
	consumer "github.com/BogoCvetkov/go_mastercalss/async/tasks/consumer"
	producer "github.com/BogoCvetkov/go_mastercalss/async/tasks/producer"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {

	// Load Config
	config := config.LoadConfig()

	// Initialize DB connection and Store
	conn, err := sql.Open(config.DBDriver, config.DB)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed connecting to DB")
	}

	store := db.NewStore(conn)

	srv := async.NewServer(store, config)

	mux := asynq.NewServeMux()

	taskManager := consumer.NewTaskConsumerManager(srv)

	mux.HandleFunc(producer.TypeEmailDelivery, taskManager.EmailProcessor.ProcessTask)

	fmt.Printf("Redis Workers start listening on Redis Queue ---> %s", config.Redis)

	if err := srv.Srv.Run(mux); err != nil {
		fmt.Printf("could not run redis workers server: %v", err)
	}
}
