package api

import (
	api "github.com/BogoCvetkov/go_mastercalss/api/controller"
	"github.com/BogoCvetkov/go_mastercalss/middleware"
)

func (s *Server) AttachRoutes() {

	ctr := api.InitControllers(s)

	api := s.router.Group("/api")

	// Public routes
	{
		public := api.Group("")

		u := public.Group("/user")
		{

			u.POST("", ctr.User.CreateUser)
			u.POST("login", ctr.User.LoginUser)
		}

	}

	// Authenticated routes
	{
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.auth, s.store))

		acc := protected.Group("/account")
		{
			acc.GET("", ctr.Account.ListAccounts)
			acc.GET("/:id", ctr.Account.GetAccount)
			acc.POST("", ctr.Account.CreateAccount)
		}

		tr := protected.Group("/transfer")
		{

			tr.POST("", ctr.Transfer.CreateTransfer)
		}

	}

}

func (s *Server) AttachGlobalMiddlewares() {
	// Handle errors globaly
	s.router.Use(middleware.ErrorMiddleware)
}
