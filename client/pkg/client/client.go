package client

import (
	"bufio"
	"io"
	"net"
	"os"
	"sync"

	"github.com/hrvadl/go-chat/server/pkg/server"
)

type Client struct {
	conn net.Conn
	wg   sync.WaitGroup
}

func NewClient() (*Client, error) {
	conn, err := net.Dial(server.Protocol, server.Port)

	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, wg: sync.WaitGroup{}}, nil
}

func (c *Client) ReceiveMessage() {
	for _, err := io.Copy(os.Stdout, c.conn); err != nil; {
	}

	c.wg.Done()
}

func (c *Client) SendMessage() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		c.conn.Write(scanner.Bytes())
	}

	c.wg.Done()
}

func (c *Client) StartMessaging() {
	c.wg.Add(2)
	go c.ReceiveMessage()
	go c.SendMessage()
	c.wg.Wait()
}
