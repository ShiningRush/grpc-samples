package server

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"gitlab.followme.com/FollowmeGo/golib/utils/network"
)

func InitPprofService() {
	serverAddr := network.GetOutboundIP()

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("server.InitPprofService:failed to listen: %v", err)
	}

	log.Println("pprof service is listen on :" + lis.Addr().String())
	go func() {
		http.Serve(lis, nil)
	}()
}
