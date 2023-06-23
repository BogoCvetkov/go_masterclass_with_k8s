package api

import (
	"github.com/BogoCvetkov/go_mastercalss/controller"
	"github.com/BogoCvetkov/go_mastercalss/middleware"
)

func (s *Server) AttachRoutes() {

	ctr := controller.InitControllers(s.Store)

	api := s.router.Group("/api")

	{
		acc := api.Group("/account")
		{
			acc.GET("", ctr.Account.ListAccounts)
			acc.GET("/:id", ctr.Account.GetAccount)
			acc.POST("", ctr.Account.CreateAccount)
		}

		tr := api.Group("/transfer")
		{

			tr.POST("", ctr.Transfer.CreateTransfer)
		}

		u := api.Group("/user")
		{

			u.POST("", ctr.User.CreateUser)
			u.POST("login", ctr.User.LoginUser)
		}
	}
}

func (s *Server) AttachGlobalMiddlewares() {
	// Handle errors globaly
	s.router.Use(middleware.ErrorMiddleware)
}
