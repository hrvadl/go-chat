package connection

import (
	"net"
	"sync"
)

type Participant struct {
	Conn    net.Conn
	Address string
}

func NewParticipant(conn net.Conn) *Participant {
	return &Participant{
		Conn:    conn,
		Address: conn.RemoteAddr().String(),
	}
}

func (p *Participant) Leave() error {
	return p.Conn.Close()
}

func (p *Participant) Send(message []byte) error {
	_, err := p.Conn.Write(message)
	return err
}

type ParticipantMap struct {
	sync.Mutex
	p map[string]*Participant
}

func NewParticipantMap() *ParticipantMap {
	return &ParticipantMap{
		p: make(map[string]*Participant),
	}
}

func (m *ParticipantMap) Store(key string, participant *Participant) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.p[key] = participant
}

func (m *ParticipantMap) Get(key string) (*Participant, bool) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	val, ok := m.p[key]
	return val, ok
}

func (m *ParticipantMap) Remove(key string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	delete(m.p, key)
}

func (m *ParticipantMap) Range(iterator func(key string, p *Participant)) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	for address, participant := range m.p {
		iterator(address, participant)
	}
}

func (m *ParticipantMap) Length() int {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	return len(m.p)
}
