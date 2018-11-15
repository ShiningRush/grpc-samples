package main

import (
	"context"
	"errors"
	"log"

	"github.com/ShiningRush/grpc-samples/pkg/config"
	fmlog "github.com/shiningrush/grpc-samples/pkg/grpc/log"
	pb "github.com/shiningrush/grpc-samples/proto/go"

	fmgrpc "github.com/shiningrush/grpc-samples/pkg/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	return &pb.HelloReply{Message: "Hello " + in.Name}, errors.New("Test")
}

func main() {
	log.Println("Test config : " + config.GetString("Test2"))
	config.InitConfig()

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			fmlog.AccessLog(
				fmlog.SetLogReqAndResp(true),
			),
		),
	}

	err := fmgrpc.Serve(
		fmgrpc.RegisterOption(opts),
		fmgrpc.RegisterService(func(s *grpc.Server) {
			pb.RegisterGreeterServer(s, &server{})
		}),
	)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
