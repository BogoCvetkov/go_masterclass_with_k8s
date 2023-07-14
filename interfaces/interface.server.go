package interfaces

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type IServer interface {
	GetStore() db.IStore
	GetAuth() IAuth
	GetConfig() *config.Config
	GetAsync() *asynq.Client
	GetRouter() *gin.Engine
}
