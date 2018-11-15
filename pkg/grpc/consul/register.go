package consul

import (
	"fmt"
	"log"
	"net"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func RegisterToConsul(consulAddr string, srvAddr string, srvName string) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul register: consul client error : ", err)
	}

	host, port, err := net.SplitHostPort(srvAddr)
	iPort, _ := strconv.Atoi(port)
	if err != nil {
		log.Fatal("consul register: service address error : ", err)
	}

	registration := &consulapi.AgentServiceRegistration{
		Name:    srvName,
		ID:      fmt.Sprintf("%v:%v", srvName, srvAddr),
		Tags:    []string{"v-" + srvName},
		Address: host,
		Port:    iPort,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%v:%v/", host, port),
			Timeout:                        "5s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "60s", // consul auto check interval is 30s, this value's is least 60s
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("consul register: register server error : ", err)
	}
}

func DeregisterFromConsul(consulAddr string, srvAddr string, srvName string) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul register : consul client : ", err)
	}

	err = client.Agent().ServiceDeregister(fmt.Sprintf("%v:%v", srvName, srvAddr))
	if err != nil {
		log.Fatal("consul register : deregister server : ", err)
	}
}
