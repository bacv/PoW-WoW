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
