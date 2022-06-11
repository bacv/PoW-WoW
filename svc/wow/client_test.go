package wow

import (
	"testing"

	"github.com/bacv/pow-wow/lib"
	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/svc/mock"
	"github.com/stretchr/testify/assert"
)

func TestClientHandler(t *testing.T) {
	w := &mock.ResponseWriter{}
	svc := NewWowClientService()

	header := hashcash.NewHashcash(mock.MockID, 1).GetHeader()
	svc.Handle(w, lib.NewChallengeMsg(header))
	mt, _, err := w.Written.Unmarshal()
	assert.NoError(t, err)
	assert.Equal(t, lib.MsgProof, mt, "message type should be MsgProof")
}
