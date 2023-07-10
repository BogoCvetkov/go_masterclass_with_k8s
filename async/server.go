package async

import (
	"context"

	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type AsyncServer struct {
	store  *db.Store
	config *config.Config
	Srv    *asynq.Server
}

// Getters that satisfy the interface
func (s *AsyncServer) GetStore() *db.Store {
	return s.store
}

func (s *AsyncServer) GetConfig() *config.Config {
	return s.config
}

func NewServer(s *db.Store, c *config.Config) *AsyncServer {

	redisAddr := c.Redis

	logger := NewLogger()
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
		})

	return &AsyncServer{
		store:  s,
		config: c,
		Srv:    srv,
	}
}
