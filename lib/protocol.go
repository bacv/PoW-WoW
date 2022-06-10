package lib

import "errors"

// MessageType ...
type MessageType byte

// Message contains powwow communication message type and it's body.
type Message []byte

const (
	// Newline representation in hex.
	// It's used to check for the end of a message that is comming through this transport.
	ByteLF = byte(0x0A)

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

	// if a message is shorter than two bytes then it has no body.
	if len(m) < 2 {
		return mt, "", mt.Validate()
	}

	// removing the new line char when unmarshaling.
	return mt, string(m[1 : len(m)-1]), mt.Validate()
}

// Marshal creates a new powwow protocol message out of message type and body.
func (m *Message) Marshal(mt MessageType, body string) error {
	buf := append([]byte{byte(mt)}, []byte(body)...)
	buf = append(buf, ByteLF)
	*m = Message(buf)
	return mt.Validate()
}

// NewMessage ...
func NewMessage(mt MessageType, body string) (Message, error) {
	m := Message{}
	err := m.Marshal(mt, body)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// NewRequestMsg ...
func NewRequestMsg() Message {
	m, _ := NewMessage(MsgRequest, "")
	return m
}

// NewChallengeMsg ...
func NewChallengeMsg(body string) Message {
	m, _ := NewMessage(MsgChallenge, body)
	return m
}

// NewProofMsg ...
func NewProofMsg(body string) Message {
	m, _ := NewMessage(MsgProof, body)
	return m
}

// NewWordsMsg ...
func NewWordsMsg(body string) Message {
	m, _ := NewMessage(MsgWords, body)
	return m
}

// NewErrorMsg ...
func NewErrorMsg(body string) Message {
	m, _ := NewMessage(MsgError, body)
	return m
}
