package interfaces

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type IGService interface {
	PassServerConfig(IGServer)
	RegisterService()
	RegisterServiceOnGateway(context.Context, *runtime.ServeMux) error
}
