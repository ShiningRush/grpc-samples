package grpc

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {

}

func TestServe(t *testing.T) {
	flag.Set("consul_addr", "192.168.8.6:8500")
	flag.Set("server_name", "User_Test")
	err := Serve()
	assert.NoError(t, err, "Serve should not get error")
}

func TestGetOutboundIP(t *testing.T) {
	ip := getOutboundIP()
	assert.NotZero(t, ip, "GetOutboundIP should not be zero value")
}
