package svc

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bacv/pow-wow/lib"
	"github.com/bacv/pow-wow/lib/hashcash"
)

type WowService interface {
	Handle(ResponseWriter, lib.Message)
}

type WisdomSource interface {
	GetWisdom() string
}

type serverSvc struct {
	mu         sync.RWMutex
	challenges map[string]*hashcash.Hashcash
	source     WisdomSource
}

func NewWowServerService(source WisdomSource) *serverSvc {
	return &serverSvc{
		challenges: make(map[string]*hashcash.Hashcash),
	}
}

func (s *serverSvc) Handle(w ResponseWriter, r lib.Message) {
	mt, m, err := r.Unmarshal()
	if err != nil {
		w.Close()
		return
	}

	switch mt {
	case lib.MsgRequest:
		header := s.addConn()
		w.Write(lib.NewChallengeMsg(header))
	case lib.MsgProof:
		if s.validate(m) {
			w.Write(lib.NewWordsMsg(s.getWisdom()))
		}

		s.removeConn(m)
		w.Close()
	default:
		w.Close()
	}
}

func (s *serverSvc) addConn() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID()
	// A very primitive way to dynamically increase the requirements for the proof.
	hash := hashcash.NewHashcash(id, uint(len(s.challenges)))
	header := hash.GetHeader()
	s.challenges[id] = hash
	return header
}

func (s *serverSvc) removeConn(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.challenges, id)
}

func (s *serverSvc) validate(m string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, err := getID(m)
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

func getID(msg string) (string, error) {
	values := strings.Split(msg, ":")
	if len(values) != 7 {
		return "", errors.New("Invalid length")
	}

	return values[3], nil
}

// TODO: append random string
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

type wisdomSource struct {
	list []string
}

func NewWisdomSource() WisdomSource {
	return &wisdomSource{
		list: []string{
			"Aim for your dreams, but don't lose yourself along the way.",
			"The best way out is always through.",
		},
	}
}

func (w *wisdomSource) GetWisdom() string {
	i := rand.Intn(len(w.list))
	if i < 1 {
		i++
	}
	return w.list[i-1]
}
