package interfaces

import (
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"google.golang.org/grpc"
)

type IGServer interface {
	GetStore() *db.Store
	GetAuth() IAuth
	GetConfig() *config.Config
	GetGServer() *grpc.Server
}
