package svc

import "github.com/bacv/pow-wow/lib"

type WowService interface {
	Handle(ResponseWriter, lib.Message)
}

type WisdomSource interface {
	GetWisdom() string
}

type IDGenerator interface {
	GenID() string
}
