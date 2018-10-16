package grpc

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/shiningrush/grpc-samples/pkg/grpc/consul"

	"github.com/shiningrush/grpc-samples/pkg/grpc/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	consulAddr string
	serverAddr string
	serverName string
)

func initParams() {
	consulAddr = utils.GetEnvOrDefault("CONSUL_ADDR", "127.0.0.1:8500")
	serverAddr = utils.GetEnvOrDefault("SERVER_ADDR", "")
	serverName = utils.GetEnvOrDefault("SERVER_NAME", "Grpc-Sample-Service")
}

type grpcServerOption struct {
	grpcOpts []grpc.ServerOption
	srvSet   SetService
}

type SetOption func(*grpcServerOption)
type SetService func(srv *grpc.Server)

func Serve(opts ...SetOption) error {
	initParams()
	srvOpt := initOpts(opts, &grpcServerOption{})
	if serverAddr == "" {
		serverAddr = getOutboundIP()
	}

	ip, _, err := net.SplitHostPort(serverAddr)
	if err != nil {
		return errors.New("grpc.Serve: server address is incorrect: " + err.Error())
	}

	s := grpc.NewServer(srvOpt.grpcOpts...)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	if srvOpt.srvSet != nil {
		srvOpt.srvSet(s)
	}

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("grpc.Serve:failed to listen: %v", err)
	}

	port := lis.Addr().(*net.TCPAddr).Port
	consul.RegisterToConsul(consulAddr, fmt.Sprintf("%v:%v", ip, port), serverName)

	log.Println(fmt.Sprintf("server is listen on : %v:%v", ip, port))
	defer func() {
		consul.DeregisterFromConsul(consulAddr, fmt.Sprintf("%v:%v", ip, port), serverName)
	}()
	return s.Serve(lis)
}

func initOpts(opts []SetOption, srvOpt *grpcServerOption) *grpcServerOption {
	for _, o := range opts {
		o(srvOpt)
	}

	return srvOpt
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	host, _, _ := net.SplitHostPort(localAddr)

	return host + ":0"
}
