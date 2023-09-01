package server

import (
	"log"
	"net"

	"github.com/hrvadl/go-chat/server/pkg/connection"
	"github.com/hrvadl/go-chat/server/pkg/messages"
)

const (
	Port     = ":5000"
	MaxConn  = 10
	Protocol = "tcp"
)

type TCP struct {
	port        string
	maxConn     int
	connections connection.ParticipantMap
	listener    net.Listener
	message     chan *messages.Message
	done        chan struct{}
}

func NewTCP(port string, maxConn int) *TCP {
	return &TCP{
		port:        port,
		maxConn:     maxConn,
		connections: *connection.NewParticipantMap(),
		message:     make(chan *messages.Message),
	}
}

func (s *TCP) Listen() {
	listener, err := net.Listen(Protocol, s.port)

	if err != nil {
		log.Fatalf("Error listening on %v: %v", s.port, err)
	}

	s.listener = listener

	for {
		participant, err := s.AcceptParticipant()

		if err != nil {
			continue
		}

		go func() {
			s.connections.Store(participant.Address, participant)
			s.ListenToMessages(participant)
		}()
	}
}

func (s *TCP) AcceptParticipant() (*connection.Participant, error) {
	conn, err := s.listener.Accept()

	if err != nil {
		return nil, err
	}

	return connection.NewParticipant(conn), nil
}

func (s *TCP) ListenToMessages(participant *connection.Participant) {
	defer participant.Leave()
	s.message <- messages.New([]byte("Greet our new member: "+participant.Address), nil, messages.Joined)

	for {
		buf := make([]byte, 32768)

		if _, err := participant.Conn.Read(buf); err != nil {
			break
		}

		s.message <- messages.New(buf, participant, messages.Default)
	}
}

func (s *TCP) BroadCast() {
	for message := range s.message {
		switch message.MsgType {
		case messages.Joined:
			fallthrough
		case messages.Left:
			s.sendToAll(message.Text)
		case messages.Default:
			s.send(message.Text, *s.exceptSender(message.Author))
		}
	}
}

func (s *TCP) sendToAll(message []byte) {
	s.connections.Range(func(key string, p *connection.Participant) {
		connection := p
		go func() {
			connection.Send(message)
		}()
	})
}

func (s *TCP) send(message []byte, receivers []connection.Participant) {
	for _, connection := range receivers {
		connection := connection
		go func() {
			connection.Send(message)
		}()
	}
}

func (s *TCP) exceptSender(sender *connection.Participant) *[]connection.Participant {
	receivers := make([]connection.Participant, 0, s.connections.Length()-1)

	s.connections.Range(func(key string, p *connection.Participant) {
		if p.Address != sender.Address {
			receivers = append(receivers, *p)
		}
	})

	return &receivers
}
