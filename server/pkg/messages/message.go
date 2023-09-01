package messages

import "github.com/hrvadl/go-chat/server/pkg/connection"

type MessageType string

const (
	Joined  MessageType = "Joined"
	Left                = "Left"
	Default             = "Default"
)

type Message struct {
	Text    []byte
	Author  *connection.Participant
	MsgType MessageType
}

func New(text []byte, author *connection.Participant, msgType MessageType) *Message {
	msg := append(text, []byte("\n")...)

	if author != nil {
		msg = append([]byte(author.Address+": "), msg...)
	}

	return &Message{
		Text:    msg,
		Author:  author,
		MsgType: msgType,
	}
}
