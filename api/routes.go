package api

import (
	"github.com/BogoCvetkov/go_mastercalss/controller"
)

func (s *Server) AttachRoutes() {

	ctr := controller.InitControllers(s.Store)

	api := s.router.Group("/api")

	{
		api.GET("/account", ctr.Account.ListAccounts)
		api.GET("/account/:id", ctr.Account.GetAccount)
		api.POST("/account", ctr.Account.CreateAccount)
	}
}
