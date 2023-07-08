package interfaces

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/hibiken/asynq"
)

type IServer interface {
	GetStore() *db.Store
	GetAuth() IAuth
	GetConfig() *config.Config
	GetAsync() *asynq.Client
}
