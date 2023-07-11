package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"

	"github.com/BogoCvetkov/go_mastercalss/auth"
	"github.com/BogoCvetkov/go_mastercalss/config"
	"github.com/BogoCvetkov/go_mastercalss/db"
	middleware "github.com/BogoCvetkov/go_mastercalss/grpc/gateway"
	interceptors "github.com/BogoCvetkov/go_mastercalss/grpc/interceptors"
	grpc_server "github.com/BogoCvetkov/go_mastercalss/grpc/services"
	"github.com/BogoCvetkov/go_mastercalss/interfaces"
	_ "github.com/BogoCvetkov/go_mastercalss/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
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
	async    *asynq.Client
}

func NewServer(s *db.Store, c *config.Config) *GRPCServer {

	a := auth.NewPasetoAuth(c.TOKEN_SECRET)
	asq := asynq.NewClient(asynq.RedisClientOpt{Addr: c.REDIS})

	return &GRPCServer{
		store:  s,
		config: c,
		auth:   a,
		async:  asq,
		services: []interfaces.IGService{
			&grpc_server.UserService{},
			&grpc_server.AccountService{},
			&grpc_server.TransferService{},
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

	// Init interceptors
	interceptManager := interceptors.InterceptorManager{}
	interceptManager.PassServerConfig(g)

	g.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptManager.NewLoggerInterceptor(), interceptManager.NewAuthInterceptor()),
	)

	// Enable reflection
	if g.config.ENV == "DEV" {
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

	// Create router
	router := mux.NewRouter()

	// Apply auth middleware to specific routes
	router.PathPrefix("/account").Handler(middleware.AuthMiddleware(gwmux, g.GetAuth(), g.GetStore()))
	router.PathPrefix("/account/{id}").Handler(middleware.AuthMiddleware(gwmux, g.GetAuth(), g.GetStore()))
	router.PathPrefix("/transfer").Handler(middleware.AuthMiddleware(gwmux, g.GetAuth(), g.GetStore()))

	// Add Swagger
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(statikFS)
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))

	// Set the gwmux as the default handler for other routes
	router.PathPrefix("/").Handler(gwmux)

	// Start Listening
	port := fmt.Sprintf(":%s", p)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot create listener")
	}

	fmt.Printf("Starting GRPC Gateway server on port --> %s \n", p)

	err = http.Serve(listener, router)
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

// Allow to send tasks to redis queue
func (g *GRPCServer) InitAsyncClient() {
	redisAddr := os.Getenv("REDIS")

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	g.async = client
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
func (s *GRPCServer) GetAsync() *asynq.Client {
	return s.async
}
