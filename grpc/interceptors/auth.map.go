package grpc_server

// Defines the auth rules for methods
var authMap = map[string]bool{
	"/UserService/CreateUser": true,
}
