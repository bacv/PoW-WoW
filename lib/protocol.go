package lib

import "errors"

// MessageType ...
type MessageType byte

// Message contains powwow communication message type and it's body.
type Message []byte

const (
	// MsgRequest is used by the client to request a challenge from the server.
	MsgRequest = MessageType(0x01)
	// MsgChallenge is sent by the server with a resource and required bit count for the proof.
	MsgChallenge = MessageType(0x02)
	// MsgProof is sent by the client with a hashcash of a received resource.
	MsgProof = MessageType(0x03)
	// MsgWords is sent by the server with the words of wisdom as a body.
	MsgWords = MessageType(0x04)
	// MsgError is a reserved message type for errors.
	MsgError = MessageType(0x05)
)

// ErrorMessageTypeUnknown ...
var ErrorMessageTypeUnknown = errors.New("Unknown message type")

// Validate checks if message type is known.
func (m MessageType) Validate() error {
	switch m {
	case MsgRequest, MsgChallenge, MsgProof, MsgWords, MsgError:
		return nil
	default:
		return ErrorMessageTypeUnknown
	}
}

// Unmarshal parses the message and returns message type with a body if message type is valid.
func (m Message) Unmarshal() (MessageType, string, error) {
	mt := MessageType(m[0])
	return mt, string(m[1:]), mt.Validate()
}

// Marshal creates a new powwow protocol message out of message type and body.
func (m *Message) Marshal(mt MessageType, body string) error {
	buf := append([]byte{byte(mt)}, []byte(body)...)
	*m = Message(buf)
	return mt.Validate()
}
