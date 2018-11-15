package consul

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.followme.com/FollowmeGo/golib/grpc/utils"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

var (
	consulAddr string
	inited     bool
	mux        sync.Mutex
)

type consulBuilder struct {
	client      *consulapi.Client
	serviceName string
}

func InitAndRegister() {
	if inited {
		return
	}

	mux.Lock()
	if !inited {
		consulAddr = utils.GetEnvOrDefault("CONSUL_ADDR", "127.0.0.1:8500")
		builder := NewConsulBuilder(consulAddr)
		resolver.Register(builder)
		inited = true
	}
	mux.Unlock()
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
	consulResolver := NewConsulResolver(&cc, target.Endpoint, opts, cb.client)
	consulResolver.wg.Add(1)
	go consulResolver.watcher()

	return consulResolver, nil
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

type consulResolver struct {
	consulClient *consulapi.Client
	clientConn   *resolver.ClientConn
	service      string
	t            *time.Timer
	rn           chan bool
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

func NewConsulResolver(cc *resolver.ClientConn, service string, opts resolver.BuildOption, consulClient *consulapi.Client) *consulResolver {
	ctx, cancel := context.WithCancel(context.Background())
	return &consulResolver{
		consulClient: consulClient,
		clientConn:   cc,
		service:      service,
		t:            time.NewTimer(time.Second * 15), // fresh frequency
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (cr *consulResolver) watcher() {
	defer cr.wg.Done()
	for {
		select {
		case <-cr.ctx.Done():
			return
		case <-cr.rn:
		case <-cr.t.C:
			cr.t.Reset(time.Second * 15)
		}
		adds, err := cr.resolve()
		if err != nil {
			log.Fatal("grpc: query service entries error:", err.Error())
		}
		(*cr.clientConn).NewAddress(adds)
		(*cr.clientConn).NewServiceConfig("")
	}
}

func (cr *consulResolver) resolve() ([]resolver.Address, error) {
	serviceEntries, _, err := cr.consulClient.Health().Service(cr.service, "", true, &consulapi.QueryOptions{})
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

func (cr *consulResolver) ResolveNow(rno resolver.ResolveNowOption) {
	select {
	case cr.rn <- true:
	default:
	}
}

func (cr *consulResolver) Close() {
	cr.cancel()
	cr.wg.Wait()
	cr.t.Stop()
}
