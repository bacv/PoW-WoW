package wow

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bacv/pow-wow/lib"
	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/svc"
)

type clientSvc struct {
}

func NewWowClientService() svc.WowService {
	return &clientSvc{}
}

func (s *clientSvc) Handle(w svc.ResponseWriter, r lib.Message) {
	mt, m, err := r.Unmarshal()

	if err != nil {
		return
	}

	switch mt {
	case lib.MsgChallenge:
		s.handleMsgChallenge(w, m)
	case lib.MsgWords:
		s.handleMsgWords(w, m)
	default:
		w.Close()
	}
}

func (s *clientSvc) handleMsgChallenge(w svc.ResponseWriter, m string) {
	values := strings.Split(m, ":")
	bits, err := strconv.ParseUint(values[1], 10, 8)
	if err != nil {
		return
	}
	header := hashcash.NewHashcash(values[3], uint(bits)).Compute()
	w.Write(lib.NewProofMsg(header))
}

func (s *clientSvc) handleMsgWords(w svc.ResponseWriter, m string) {
	fmt.Println(m)
}
