package main

import (
	"log"

	"github.com/hrvadl/go-chat/server/pkg/server"
)

const (
	port    = ":5000"
	maxConn = 10
)

func main() {
	tcpServer := server.NewTCP(port, maxConn)

	if err := tcpServer.Listen(); err != nil {
		log.Fatalf("Error listening on %v: %v", port, err)
		return
	}

}
