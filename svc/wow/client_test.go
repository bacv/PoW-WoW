package wow

import (
	"crypto/sha1"
	"testing"

	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/lib/protocol"
	"github.com/bacv/pow-wow/svc/mock"
	"github.com/stretchr/testify/assert"
)

func TestClientHandler(t *testing.T) {
	w := &mock.ResponseWriter{}
	svc := NewWowClientService()

	header := hashcash.NewHashcash(mock.MockID, 1, sha1.New()).GetHeader()
	svc.Handle(w, protocol.NewChallengeMsg(header))
	mt, _, err := w.Written.Unmarshal()
	assert.NoError(t, err)
	assert.Equal(t, protocol.MsgProof, mt, "message type should be MsgProof")
}
