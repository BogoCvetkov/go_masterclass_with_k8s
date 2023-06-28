package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/BogoCvetkov/go_mastercalss/auth"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	grpc_server "github.com/BogoCvetkov/go_mastercalss/grpc/services"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type GRPCServer struct {
	store    *db.Store
	auth     interfaces.IAuth
	config   *config.Config
	srv      *grpc.Server
	services []interfaces.IGService
}

func NewServer(s *db.Store, c *config.Config) *GRPCServer {

	a := auth.NewPasetoAuth(c.TokenSecret)

	return &GRPCServer{
		store:  s,
		config: c,
		auth:   a,
		services: []interfaces.IGService{
			&grpc_server.UserService{},
		},
	}
}

// GRPC Server
func (g *GRPCServer) Start(p string) {
	fmt.Printf("Starting GRPC server on port --> %s \n", p)

	port := fmt.Sprintf(":%s", p)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Panic(err)
	}

	g.srv = grpc.NewServer()
	// Enable reflection
	if g.config.Env == "DEV" {
		reflection.Register(g.srv)
	}

	g.initServices()

	if err := g.srv.Serve(lis); err != nil {
		log.Fatalf("grpc failed to serve: %v", err)
	}

}

// HTTP Gateway GRPC Server
func (g *GRPCServer) StartHttpGateway(p string) {

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	gwmux := runtime.NewServeMux(jsonOption)

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := g.initGatewayServices(c, gwmux)
	if err != nil {
		log.Fatalf("cannot register Gateway handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	port := fmt.Sprintf(":%s", p)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create listener")
	}

	fmt.Printf("Starting GRPC Gateway server on port --> %s \n", p)

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start GRPC Gateway")
	}

}

// Register Service Methods with the GRPC Server
func (g *GRPCServer) initServices() {

	for _, s := range g.services {
		s.PassServerConfig(g)
		s.RegisterService()
	}
}

func (g *GRPCServer) initGatewayServices(c context.Context, mux *runtime.ServeMux) error {

	for _, s := range g.services {
		s.PassServerConfig(g)
		if err := s.RegisterServiceOnGateway(c, mux); err != nil {
			return err
		}
	}

	return nil
}

// Getters that satisfy the interface
func (s *GRPCServer) GetStore() *db.Store {
	return s.store
}
func (s *GRPCServer) GetAuth() interfaces.IAuth {
	return s.auth
}
func (s *GRPCServer) GetConfig() *config.Config {
	return s.config
}
func (s *GRPCServer) GetGServer() *grpc.Server {
	return s.srv
}
