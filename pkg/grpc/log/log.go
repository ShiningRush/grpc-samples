package log

import (
	"context"
	"github.com/shiningrush/grpc-samples/pkg/log/access"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

type accessLogOption struct {
	enableReqAndResp bool
}

type setOption func(*accessLogOption)

func SetLogReqAndResp(v bool) setOption {
	return func(o *accessLogOption) {
		o.enableReqAndResp = v
	}
}

func AccessLog(setOpts ...setOption) grpc.UnaryServerInterceptor {
	opts := initOption(setOpts)

	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && fullMethodName == "/grpc.health.v1.Health/Check" {
				return false
			}

			return true
		}),
	}

	// grpc_zap.ReplaceGrpcLogger(access.Logger)
	return grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(access.Logger, zapOpts...),
		grpc_zap.PayloadUnaryServerInterceptor(access.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			if !opts.enableReqAndResp {
				return false
			}

			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if fullMethodName == "/grpc.health.v1.Health/Check" {
				return false
			}

			return true
		}),
	)
}

func initOption(setOpts []setOption) *accessLogOption {
	opts := &accessLogOption{}
	for _, v := range setOpts {
		v(opts)
	}

	return opts
}
