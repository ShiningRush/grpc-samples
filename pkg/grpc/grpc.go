package grpc

import (
	"errors"
	"fmt"
	"log"
	"net"

	"gitlab.followme.com/FollowmeGo/golib/grpc/consul"
	"gitlab.followme.com/FollowmeGo/golib/utils/network"

	"gitlab.followme.com/FollowmeGo/golib/grpc/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	consulAddr string
	serverAddr string
)

func initParams() {
	consulAddr = utils.GetEnvOrDefault("CONSUL_ADDR", "127.0.0.1:8500")
	serverAddr = utils.GetEnvOrDefault("SERVER_ADDR", "")
}

type grpcServerOption struct {
	name     string
	grpcOpts []grpc.ServerOption
	srvSet   SetService
}

type SetOption func(*grpcServerOption)
type SetService func(srv *grpc.Server)

func Serve(opts ...SetOption) error {
	initParams()
	srvOpt := initOpts(opts, &grpcServerOption{name: "Grpc-sample-srv"})
	if serverAddr == "" {
		serverAddr = network.GetOutboundIP()
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
	consul.RegisterToConsul(consulAddr, fmt.Sprintf("%v:%v", ip, port), srvOpt.name)

	log.Println(fmt.Sprintf("server is listen on : %v:%v", ip, port))
	defer func() {
		consul.DeregisterFromConsul(consulAddr, fmt.Sprintf("%v:%v", ip, port), srvOpt.name)
	}()
	return s.Serve(lis)
}

func initOpts(opts []SetOption, srvOpt *grpcServerOption) *grpcServerOption {
	for _, o := range opts {
		o(srvOpt)
	}

	return srvOpt
}

func Dial(srvName string) *grpc.ClientConn {
	consul.InitAndRegister()

	conn, err := grpc.Dial("consul:///"+srvName, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("Dial srv: %v failed : %v", srvName, err.Error())
	}

	return conn
}
