package hashcash

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	resource           string
	bits               uint
	expectedZeros      int
	lastCharIdx        int
	acceptableEndChars string
}

var commonCases = []testCase{
	{"test1", 0, 0, 0, "01"},
	{"test1", 1, 0, 0, "01"},
	{"test2", 2, 0, 0, "0123"},
	{"test3", 3, 0, 0, "01234567"},
	{"test4", 4, 1, 0, "0"},
	{"test5", 11, 2, 2, "01234567"},
	{"test6", 13, 2, 2, "01"},
	{"test7", 14, 2, 2, "0123"},
}

func TestGetHeader(t *testing.T) {
	tc := commonCases[0]
	hash := NewHashcash(tc.resource, tc.bits)
	emptyHeader := hash.GetHeader()
	proofHeader := hash.Compute()

	for _, header := range []string{emptyHeader, proofHeader} {
		values := strings.Split(header, ":")
		assert.Equal(t, len(values), 7, "header should contain 7 fields seperated by ':'")
		assert.Equal(t, values[1], fmt.Sprintf("%d", tc.bits), fmt.Sprintf("%s should be %d", values[1], tc.bits))
		assert.Equal(t, values[3], tc.resource, fmt.Sprintf("%s should be %s", values[3], tc.resource))
	}
}

func TestCheckZeros(t *testing.T) {
	for _, tc := range commonCases {
		hash := NewHashcash(tc.resource, tc.bits)
		proof := hash.Compute()
		digest := hash.Digest(proof)

		// check expected zeros
		if tc.expectedZeros > 0 {
			zeros := "0"
			for i := 1; i < tc.expectedZeros; i++ {
				zeros = zeros + "0"
			}
			assert.True(t, digest[:tc.expectedZeros] == zeros, fmt.Sprintf("header should have %s as a suffix", zeros))
		}

		// check last char
		lastChar := string(digest[tc.lastCharIdx])
		assert.True(t, strings.Contains(tc.acceptableEndChars, lastChar),
			fmt.Sprintf("end char is '%s' but should be one of %s", lastChar, tc.acceptableEndChars))
	}
}
