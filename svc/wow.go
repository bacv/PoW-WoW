package svc

import "github.com/bacv/pow-wow/lib/protocol"

type WowService interface {
	Handle(ResponseWriter, protocol.Message)
}

type WisdomSource interface {
	GetWisdom() string
}

type IDGenerator interface {
	GenID() string
}
