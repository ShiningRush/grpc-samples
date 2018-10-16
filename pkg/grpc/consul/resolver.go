package consul

import (
	"context"
	"github.com/shiningrush/grpc-samples/pkg/grpc/utils"
	"fmt"
	"log"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

var (
	consulAddr string
)

type consulBuilder struct {
	client      *consulapi.Client
	serviceName string
}

func InitAndRegister() {
	consulAddr = utils.GetEnvOrDefault("CONSUL_ADDR", "127.0.0.1:8500")
	builder := NewConsulBuilder(consulAddr)
	resolver.Register(builder)
}

func NewConsulBuilder(address string) resolver.Builder {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("grpc resolver : create consul client error", err.Error())
		return nil
	}

	return &consulBuilder{client: client}
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	cb.serviceName = target.Endpoint

	consulResolver := NewConsulResolver(&cc, cb, opts)
	go consulResolver.watcher()

	return consulResolver, nil
}

func (cb consulBuilder) resolve() ([]resolver.Address, error) {

	serviceEntries, _, err := cb.client.Health().Service(cb.serviceName, "", true, &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	adds := make([]resolver.Address, 0)
	for _, serviceEntry := range serviceEntries {
		address := resolver.Address{Addr: fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port)}
		adds = append(adds, address)
	}
	return adds, nil
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

type consulResolver struct {
	clientConn    *resolver.ClientConn
	consulBuilder *consulBuilder
	t             *time.Timer
	rn            chan bool
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewConsulResolver(cc *resolver.ClientConn, cb *consulBuilder, opts resolver.BuildOption) *consulResolver {
	ctx, cancel := context.WithCancel(context.Background())
	return &consulResolver{
		clientConn:    cc,
		consulBuilder: cb,
		t:             time.NewTimer(0),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (cr *consulResolver) watcher() {
	for {
		select {
		case <-cr.ctx.Done():
			return
		case <-cr.rn:
		case <-cr.t.C:
		}
		adds, err := cr.consulBuilder.resolve()
		if err != nil {
			log.Fatal("grpc: query service entries error:", err.Error())
		}
		(*cr.clientConn).NewAddress(adds)
		(*cr.clientConn).NewServiceConfig("")
	}
}

func (cr *consulResolver) Scheme() string {
	return cr.consulBuilder.Scheme()
}

func (cr *consulResolver) ResolveNow(rno resolver.ResolveNowOption) {
	select {
	case cr.rn <- true:
	default:
	}
}

func (cr *consulResolver) Close() {
	cr.cancel()
	cr.t.Stop()
}
