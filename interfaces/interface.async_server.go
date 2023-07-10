package interfaces

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
)

type IAsyncServer interface {
	GetStore() *db.Store
	GetConfig() *config.Config
}
