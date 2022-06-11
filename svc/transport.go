package svc

import (
	"bufio"
	"errors"
	"net"
	"sync"

	"github.com/bacv/pow-wow/lib/protocol"
)

var (
	ErrorSendChannelClosed = errors.New("Send channel is closed")
	ErrorWriteToClosed     = errors.New("Writing to a closed transport")
)

type ResponseWriter interface {
	Write(protocol.Message) error
	Close()
}

type HandleFunc func(ResponseWriter, protocol.Message)

type Transport struct {
	conn      net.Conn
	sendC     chan []byte
	stopC     chan struct{}
	handler   HandleFunc
	closeOnce sync.Once
	closed    bool
	mu        sync.RWMutex
}

func NewTransport(conn net.Conn, handler HandleFunc) *Transport {
	return &Transport{
		conn:    conn,
		handler: handler,
		sendC:   make(chan []byte),
		stopC:   make(chan struct{}),
	}
}

func (t *Transport) Spawn() error {
	defer t.Close()
	errC := make(chan error)

	go func() {
		t.read(errC)
	}()

	err := <-errC
	return err
}

func (t *Transport) Close() {
	t.closeOnce.Do(func() {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.closed = true
		close(t.stopC)
		t.conn.Close()
	})
}

func (t *Transport) Write(msg protocol.Message) error {
	if t.IsClosed() {
		return ErrorWriteToClosed
	}

	t.conn.Write(msg)
	return nil
}

func (t *Transport) IsClosed() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.closed
}

func (t *Transport) read(errC chan<- error) {
	for {
		select {
		case <-t.stopC:
			return
		default:
			bytes, err := bufio.NewReader(t.conn).ReadBytes(protocol.ByteLF)
			if err != nil {
				errC <- err
				return
			}

			t.handler(t, protocol.Message(bytes))
		}
	}
}
