package main

import (
	"github.com/hrvadl/go-chat/server/pkg/server"
)

func main() {
	tcpServer := server.NewTCP(server.Port, server.MaxConn)

	go tcpServer.Listen()
	tcpServer.BroadCast()
}
