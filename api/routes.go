package api

import (
	"net/http"

	api "github.com/BogoCvetkov/go_mastercalss/api/controller"
	"github.com/BogoCvetkov/go_mastercalss/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) AttachRoutes() {

	ctr := api.InitControllers(s)

	api := s.router.Group("/api")

	// health check
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	// Public routes
	{
		public := api.Group("")

		u := public.Group("/user")
		{

			u.POST("", ctr.User.CreateUser)
			u.POST("login", ctr.User.LoginUser)
			u.POST("refresh-token", ctr.User.RefreshToken)
			u.POST("verify", ctr.User.VerifyEmail)
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
