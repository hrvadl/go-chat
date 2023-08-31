package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
)

const protocol = "tcp"

type Connection struct {
	net.Conn
	nickname string
}

func NewTCP(port string, maxConn int) *TCP {
	return &TCP{
		port:        port,
		maxConn:     maxConn,
		connections: map[string]Connection{},
	}
}

type TCP struct {
	port        string
	maxConn     int
	connections map[string]Connection
	listener    net.Listener
}

func (s *TCP) Listen() error {
	listener, err := net.Listen(protocol, s.port)

	if err != nil {
		return err
	}

	s.listener = listener

	for {
		conn, err := s.AcceptConnection()

		go func() {
			defer conn.Close()
			s.connections[conn.nickname] = *conn

			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}

			s.StartChatting(conn)
		}()
	}
}

func (s *TCP) AcceptConnection() (*Connection, error) {
	conn, err := s.listener.Accept()

	if err != nil {
		log.Fatalf("error accepting connection %v \n", err)
	}

	conn.Write([]byte("Hi, please write your nickname. Note: it should be unique\n"))

	var nickname bytes.Buffer

	for {
		if nickname, err := io.ReadAll(conn); err != nil {
			conn.Write([]byte(nickname))
			break
		}
	}

	if err := s.validateNickname(nickname.String()); err != nil {
		conn.Write([]byte(err.Error()))
		return nil, err
	}

	newConn := &Connection{
		nickname: nickname.String(),
		Conn:     conn,
	}

	return newConn, nil

}

func (s *TCP) StartChatting(conn *Connection) {
	s.sendMessage([]byte("Greeting to our new member - "+conn.nickname), conn)
	for {
		var message []byte

		if _, err := conn.Read(message); err != nil {
			conn.Write([]byte("Couldn't read message. Please try again later\n"))
			return
		}

		s.sendMessage(message, conn)
	}
}

func (s *TCP) validateNickname(nickname string) error {
	for _, c := range s.connections {
		if c.nickname == nickname {
			return errors.New("Such nickname already exists\n")
		}
	}

	return nil
}

func (s *TCP) sendMessage(message []byte, sender *Connection) {
	for _, c := range s.connections {
		c := c
		go func() {
			if c.nickname != sender.nickname {
				c.Write(message)
			}
		}()
	}
}
