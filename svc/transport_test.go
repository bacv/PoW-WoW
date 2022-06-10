package svc

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransportWriteToClosed(t *testing.T) {
	conn, _ := net.Pipe()

	transport := NewTransport(conn, func(w ResponseWriter, r Request) {})
	go func() {
		transport.Spawn()
	}()
	transport.Close()
	conn.Close()

	err := transport.Write("test")
	assert.ErrorIs(t, err, ErrorWriteToClosed, "it should not be possible to write to a clossed transport")
}
