package hashcash

import (
	"encoding/base64"
	"fmt"
	"hash"
	"math/rand"
	"strconv"
)

type Hashcash struct {
	hasher   hash.Hash
	bits     uint
	resource string
	rand     string
	counter  uint64
}

func NewHashcash(resource string, bits uint, hasher hash.Hash) *Hashcash {
	rb, _ := randomBytes(8)
	return &Hashcash{
		hasher:   hasher,
		bits:     bits,
		resource: resource,
		rand:     base64EncodeBytes(rb),
		counter:  rand.Uint64(),
	}
}

func (h *Hashcash) Compute() string {
	header := h.GetHeader()
	for !h.verifyProof(header) {
		h.counter++
		header = h.GetHeader()
	}
	return header
}

func (h *Hashcash) Verify(header string) (bool, error) {
	// TODO: check header formatting and other fields.
	return h.verifyProof(header), nil
}

func (h *Hashcash) verifyProof(header string) bool {
	digest := h.Digest(header)
	return checkZeros(digest, h.bits)
}

func (h *Hashcash) Digest(header string) string {
	h.hasher.Reset()
	h.hasher.Write([]byte(header))
	return fmt.Sprintf("%x", h.hasher.Sum(nil))
}

func (h *Hashcash) GetHeader() string {
	return fmt.Sprintf("1:%d::%s::%s:%s",
		h.bits,
		h.resource,
		h.rand,
		base64EncodeUint(h.counter),
	)
}

func checkZeros(digest string, bits uint) bool {
	// Do we need additional check for the last rune representing last four bits?
	mod := bits % 4
	runes := bits / 4
	if mod != 0 || bits == 0 {
		runes++
	}

	for _, v := range digest[:runes-1] {
		if v != '0' {
			return false
		}
	}

	lastChar := int(digest[runes-1] - '0')
	if mod == 0 && lastChar != 0 {
		return false
	} else if mod == 3 && lastChar > 7 {
		return false
	} else if mod == 2 && lastChar > 3 {
		return false
	} else if mod == 1 && lastChar > 1 {
		return false
	}

	return true
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func base64EncodeBytes(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func base64EncodeUint(n uint64) string {
	return base64EncodeBytes([]byte(strconv.FormatUint(n, 10)))
}
