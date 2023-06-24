package interfaces

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
)

type IServer interface {
	GetStore() *db.Store
	GetAuth() IAuth
	GetConfig() *config.Config
}
