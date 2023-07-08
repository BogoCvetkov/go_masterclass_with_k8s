package main

import (
	"context"
	"fmt"

	"github.com/BogoCvetkov/go_mastercalss/async"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {

	// Load Config
	config := config.LoadConfig()

	redisAddr := config.Redis

	logger := async.NewLogger()
	redis.SetLogger(logger)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 3,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("process task failed")
			}),
			Logger: logger,
		},
	)

	mux := asynq.NewServeMux()

	asyncManager := async.NewAsyncManager(config)

	mux.HandleFunc(async.TypeEmailDelivery, asyncManager.EmailProcessor.ProcessTask)

	fmt.Printf("Redis Workers start listening on Redis Queue ---> %s", redisAddr)

	if err := srv.Run(mux); err != nil {
		fmt.Printf("could not run redis workers server: %v", err)
	}
}
