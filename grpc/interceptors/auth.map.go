package grpc_server

// Defines the auth rules for methods
var authMap = map[string]bool{
	// Account
	"/AccountService/CreateAccount": true,
	"/AccountService/GetAccount":    true,
	"/AccountService/ListAccounts":  true,

	// Transfer
	"/TransferService/CreateTransfer": true,
}
