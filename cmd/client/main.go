package main

import (
	"log"
	"os"
	"time"

	pb "github.com/shiningrush/grpc-samples/proto/go"

	"github.com/shiningrush/grpc-samples/pkg/grpc/consul"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

const (
	defaultName = "world"
)

func main() {
	os.Setenv("CONSUL_ADDR", "192.168.8.6:8500")
	consul.InitAndRegister()
	// Set up a connection to the server.
	conn, err := grpc.Dial("consul:///followme-srv-user", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rsp, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", rsp.Message)
}
