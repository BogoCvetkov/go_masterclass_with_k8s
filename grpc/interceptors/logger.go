package grpc_server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *InterceptorManager) NewLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		result, err := handler(c, req)
		duration := time.Since(startTime)

		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		logger := log.Info()
		if err != nil {
			logger = log.Error().Err(err)
		}

		logger.Str("protocol", "grpc").
			Str("method", info.FullMethod).
			Int("status_code", int(statusCode)).
			Str("status_text", statusCode.String()).
			Dur("duration", duration).
			Msg("received a gRPC request")

		return result, err

	}
}
