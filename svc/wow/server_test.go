package wow

import (
	"strconv"
	"strings"
	"testing"

	"github.com/bacv/pow-wow/lib"
	"github.com/bacv/pow-wow/lib/hashcash"
	"github.com/bacv/pow-wow/svc/mock"
	"github.com/stretchr/testify/assert"
)

func TestServerHandler(t *testing.T) {
	source := &mock.Source{Words: "mocks of wisdom"}
	w := &mock.ResponseWriter{}
	svc := newWowServerService(source, &mock.MockGenerator{})

	// Imitate first request from client
	svc.Handle(w, lib.NewRequestMsg())
	_, m, err := w.Written.Unmarshal()
	assert.NoError(t, err)
	values := strings.Split(m, ":")
	assert.Equal(t, mock.MockID, values[3], "")

	// Imitate second request from client
	bits, err := strconv.ParseUint(values[1], 10, 8)
	assert.NoError(t, err)
	header := hashcash.NewHashcash(values[3], uint(bits)).Compute()
	svc.Handle(w, lib.NewProofMsg(header))
	_, m, err = w.Written.Unmarshal()
	assert.NoError(t, err)
	assert.Equal(t, source.Words, m, "")
}
