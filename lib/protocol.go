package lib

type MessageType byte
type Message []byte

const (
	MsgRequest   = MessageType(0x01)
	MsgChallenge = MessageType(0x02)
	MsgProof     = MessageType(0x03)
	MsgWords     = MessageType(0x04)
	MsgError     = MessageType(0x05)
)

type ProtocolParser interface {
	GetContent(Message) (string, error)
	GetType(Message) (MessageType, error)
}
