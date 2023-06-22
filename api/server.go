package api

import (
	"fmt"
	"log"

	"github.com/BogoCvetkov/go_mastercalss/db"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
}

type Server struct {
	Store  *db.Store
	router *gin.Engine
}

func NewServer(s *db.Store) *Server {

	r := gin.Default()

	return &Server{
		Store:  s,
		router: r,
	}
}

func (s *Server) Start(p string) {
	fmt.Printf("Starting server on port --> %s", p)

	port := fmt.Sprintf(":%s", p)

	if err := s.router.Run(port); err != nil {
		log.Panic(err)
	}
}
