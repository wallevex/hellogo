package grpc_logging

import (
	"context"
	"hellogo/pkg/x/log"

	"google.golang.org/grpc"
)

func PayloadUnaryServerInterceptor(log log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Debugw("server request infos",
			"method", info.FullMethod,
			"payload", req,
		)
		return handler(ctx, req)
	}
}
