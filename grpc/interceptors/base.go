package grpc_server

import "github.com/BogoCvetkov/go_mastercalss/interfaces"

// Used to inject all necessary dependencies
type InterceptorManager struct {
	server interfaces.IGServer
}

func (m *InterceptorManager) PassServerConfig(server interfaces.IGServer) {
	m.server = server
}
