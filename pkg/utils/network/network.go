package network

import (
	"log"
	"net"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	host, _, err := net.SplitHostPort(localAddr)
	if err != nil {
		log.Fatalf("get outbound ip failed (localAddr: %v) : %v", localAddr, err.Error())
	}
	log.Println("outbound ip : " + host)
	return host + ":0"
}
