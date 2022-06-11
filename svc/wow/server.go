package wow

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/lib/protocol"
	"github.com/bacv/pow-wow/svc"
)

type serverSvc struct {
	mu         sync.RWMutex
	challenges map[string]*hashcash.Hashcash
	generator  svc.IDGenerator
	source     svc.WisdomSource
	balancer   svc.LoadBalancer
}

type generator struct{}

func (g *generator) GenID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type balancer struct{}

func (b *balancer) GetChallengeBits(load int) uint {
	return uint(10 + load)
}

func NewWowServerService(source svc.WisdomSource) *serverSvc {
	return newWowServerService(source, &generator{}, &balancer{})
}

func newWowServerService(source svc.WisdomSource, g svc.IDGenerator, b svc.LoadBalancer) *serverSvc {
	return &serverSvc{
		challenges: make(map[string]*hashcash.Hashcash),
		generator:  g,
		source:     source,
		balancer:   b,
	}
}

func (s *serverSvc) Handle(w svc.ResponseWriter, r protocol.Message) {
	mt, m, err := r.Unmarshal()
	if err != nil {
		w.Close()
		return
	}

	switch mt {
	case protocol.MsgRequest:
		s.handleMsgRequest(w)
	case protocol.MsgProof:
		s.handleMsgProof(w, m)
	default:
		w.Close()
	}
}

func (s *serverSvc) handleMsgRequest(w svc.ResponseWriter) {
	header := s.addConn()
	w.Write(protocol.NewChallengeMsg(header))
}

func (s *serverSvc) handleMsgProof(w svc.ResponseWriter, m string) {
	if s.validate(m) {
		msg := protocol.NewWordsMsg(s.getWisdom())
		w.Write(msg)
	}

	s.removeConn(m)
	w.Close()
}

func (s *serverSvc) addConn() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.generator.GenID()
	// A very primitive way to dynamically increase the requirements for the proof.
	bits := s.balancer.GetChallengeBits(len(s.challenges))
	hash := hashcash.NewHashcash(id, bits)
	header := hash.GetHeader()
	s.challenges[id] = hash
	return header
}

func (s *serverSvc) removeConn(m string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := extractID(m)
	if err != nil {
		return
	}
	delete(s.challenges, id)
}

func (s *serverSvc) validate(m string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, err := extractID(m)
	if err != nil {
		return false
	}

	if h, ok := s.challenges[id]; ok {
		verified, err := h.Verify(m)
		if verified && err == nil {
			return true
		}
	}
	return false
}

func (s *serverSvc) getWisdom() string {
	return s.source.GetWisdom()
}

func extractID(msg string) (string, error) {
	values := strings.Split(msg, ":")
	if len(values) != 7 {
		return "", errors.New("Invalid length")
	}

	return values[3], nil
}
