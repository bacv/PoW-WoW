package protocol

import (
	"testing"

	"github.com/bacv/pow-wow/lib"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectedType    lib.MessageType
	expectedContent string
	rawMessage      []byte
}

var commonCases = []testCase{
	{lib.MsgRequest, "", []byte{0x01}},
	{lib.MsgChallenge, "test:2", []byte{0x02}},
	{lib.MsgProof, "1::::::", []byte{0x03}},
	{lib.MsgWords, "of wisdom", []byte{0x04}},
	{lib.MsgError, "some error", []byte{0x05}},
}

func TestMessageType(t *testing.T) {
	parser := NewParser()
	for _, tc := range commonCases {
		mt, err := parser.GetType(tc.rawMessage)
		assert.NoError(t, err)

		m, err := parser.GetContent(tc.rawMessage)
		assert.NoError(t, err)

		assert.Equal(t, mt, tc.expectedType, "")
		assert.Equal(t, m, tc.expectedContent, "")
	}
}
