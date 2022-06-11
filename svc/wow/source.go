package wow

import (
	"math/rand"

	"github.com/bacv/pow-wow/svc"
)

type wisdomSource struct {
	list []string
}

func NewWisdomSource() svc.WisdomSource {
	return &wisdomSource{
		list: []string{
			"Aim for your dreams, but don't lose yourself along the way.",
			"The best way out is always through.",
			"When your enemy is making a very serious mistake, don't be impolite anddisturb him.",
		},
	}
}

func (w *wisdomSource) GetWisdom() string {
	i := rand.Intn(len(w.list))
	return w.list[i]
}
