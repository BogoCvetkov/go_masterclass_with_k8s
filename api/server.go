package api

import (
	"fmt"
	"log"

	"github.com/BogoCvetkov/go_mastercalss/auth"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

type Server struct {
	router *gin.Engine
	store  *db.Store
	auth   interfaces.IAuth
	config *config.Config
	async  *asynq.Client
}

func NewServer(s *db.Store, c *config.Config) *Server {

	r := gin.Default()
	a := auth.NewPasetoAuth(c.TOKEN_SECRET)
	asq := asynq.NewClient(asynq.RedisClientOpt{Addr: c.REDIS})

	return &Server{
		store:  s,
		router: r,
		config: c,
		auth:   a,
		async:  asq,
	}
}

func (s *Server) Start(p string) {
	fmt.Printf("Starting server on port --> %s \n", p)

	port := fmt.Sprintf(":%s", p)

	if err := s.router.Run(port); err != nil {
		log.Panic(err)
	}
}

// Getters that satisfy the interface
func (s *Server) GetStore() *db.Store {
	return s.store
}
func (s *Server) GetAuth() interfaces.IAuth {
	return s.auth
}
func (s *Server) GetConfig() *config.Config {
	return s.config
}
func (s *Server) GetAsync() *asynq.Client {
	return s.async
}
