package main

import (
	"log"

	"github.com/hrvadl/go-chat/client/pkg/client"
	"github.com/hrvadl/go-chat/server/pkg/server"
)

func main() {
	client, err := client.NewClient()

	if err != nil {
		log.Fatalf("cannot connect to the server on port %v: %v", server.Port, err)
	}

	client.StartMessaging()
}
