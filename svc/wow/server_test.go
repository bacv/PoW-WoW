package wow

import (
	"crypto/sha1"
	"strconv"
	"strings"
	"testing"

	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/lib/protocol"
	"github.com/bacv/pow-wow/svc/mock"
	"github.com/stretchr/testify/assert"
)

func TestServerHandler(t *testing.T) {
	source := &mock.Source{Words: "mocks of wisdom"}
	w := &mock.ResponseWriter{}
	svc := newWowServerService(source, &mock.MockGenerator{}, &mock.MockBalancer{})

	// Imitate first request from client
	svc.Handle(w, protocol.NewRequestMsg())
	_, m, err := w.Written.Unmarshal()
	assert.NoError(t, err)
	values := strings.Split(m, ":")
	assert.Equal(t, mock.MockID, values[3], "")

	// Imitate second request from client
	bits, err := strconv.ParseUint(values[1], 10, 8)
	assert.NoError(t, err)
	header := hashcash.NewHashcash(values[3], uint(bits), sha1.New()).Compute()
	svc.Handle(w, protocol.NewProofMsg(header))
	_, m, err = w.Written.Unmarshal()
	assert.NoError(t, err)
	assert.Equal(t, source.Words, m, "")
}
