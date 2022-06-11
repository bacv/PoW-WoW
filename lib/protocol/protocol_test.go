package protocol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	expectedType    MessageType
	expectedContent string
	rawMessage      Message
}

var commonCases = []testCase{
	{MsgRequest, "", Message{0x01, ByteLF}},
	{MsgChallenge, "test:2", Message{0x02, 0x74, 0x65, 0x73, 0x74, 0x3A, 0x32, ByteLF}},
	{MsgProof, "1::::::", Message{0x03, 0x31, 0x3a, 0x3A, 0x3A, 0x3A, 0x3A, 0x3A, ByteLF}},
	{MsgWords, "of wisdom", Message{0x04, 0x6F, 0x66, 0x20, 0x77, 0x69, 0x73, 0x64, 0x6F, 0x6D, ByteLF}},
	{MsgError, "some error", Message{0x05, 0x73, 0x6F, 0x6D, 0x65, 0x20, 0x65, 0x72, 0x72, 0x6F, 0x72, ByteLF}},
}

func TestMessageUnmarshal(t *testing.T) {
	for _, tc := range commonCases {
		mt, m, err := tc.rawMessage.Unmarshal()
		assert.NoError(t, err)

		assert.Equal(t, tc.expectedType, mt, "wrong message type")
		assert.Equal(t, tc.expectedContent, m, "wrong message body")
	}
}

func TestMessageMarshal(t *testing.T) {
	for _, tc := range commonCases {
		m := Message{}
		err := m.Marshal(tc.expectedType, tc.expectedContent)

		assert.NoError(t, err)
		assert.Equal(t, tc.rawMessage, m, fmt.Sprintf("got %s, should be %s", m, tc.rawMessage))
	}
}

func TestMessageInvalidType(t *testing.T) {
	m := Message{}
	err := m.Marshal(MessageType(0x0A), "wrong type")

	assert.ErrorIs(t, err, ErrorMessageTypeUnknown, "error should be ErrorMessageTypeUnknown")
}

func TestMessageShorthands(t *testing.T) {
	for _, tc := range commonCases {
		m, err := NewMessage(tc.expectedType, tc.expectedContent)
		assert.NoError(t, err)
		assert.Equal(t, tc.rawMessage, m, "")
	}
}
