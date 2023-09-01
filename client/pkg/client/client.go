package client

import (
	"bufio"
	"io"
	"net"
	"os"

	"github.com/hrvadl/go-chat/server/pkg/server"
)

type Client struct {
	conn net.Conn
	Done chan struct{}
}

func NewClient() (*Client, error) {
	conn, err := net.Dial(server.Protocol, server.Port)

	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, Done: make(chan struct{})}, nil
}

func (c *Client) ReceiveMessage() {
	for _, err := io.Copy(os.Stdout, c.conn); err != nil; {

	}

	c.Leave()
}

func (c *Client) SendMessage() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		c.conn.Write(scanner.Bytes())
	}
}

func (c *Client) Leave() {
	c.Done <- struct{}{}
}
