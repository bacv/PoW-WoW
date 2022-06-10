package protocol

import "github.com/bacv/pow-wow/lib"

type wowParser struct{}

func NewParser() lib.ProtocolParser {
	return &wowParser{}
}

func (p *wowParser) GetContent(msg lib.Message) (string, error) {
	return "", nil
}

func (p *wowParser) GetType(msg lib.Message) (lib.MessageType, error) {
	return lib.MsgChallenge, nil
}
