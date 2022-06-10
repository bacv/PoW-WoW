package svc

import (
	"bufio"
	"errors"
	"net"
	"sync"

	"github.com/bacv/pow-wow/lib"
)

var (
	ErrorSendChannelClosed = errors.New("Send channel is closed")
	ErrorWriteToClosed     = errors.New("Writing to a closed transport")
)

type ResponseWriter interface {
	Write(lib.Message) error
	Close()
}

type HandleFunc func(ResponseWriter, lib.Message)

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
	var wg sync.WaitGroup
	wg.Add(2)
	errC := make(chan error, 2)

	go func() {
		defer wg.Done()
		t.read(errC)
	}()

	go func() {
		defer wg.Done()
		t.write(errC)
	}()

	go func() {
		wg.Wait()
		close(errC)
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

func (t *Transport) Write(msg lib.Message) error {
	if t.IsClosed() {
		return ErrorWriteToClosed
	}

	t.sendC <- []byte(msg)
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
			bytes, err := bufio.NewReader(t.conn).ReadBytes(lib.ByteLF)
			if err != nil {
				errC <- err
				return
			}

			t.handler(t, lib.Message(bytes))
		}
	}
}

func (t *Transport) write(errC chan<- error) {
	for {
		select {
		case msg, ok := <-t.sendC:
			if !ok {
				errC <- ErrorSendChannelClosed
				return
			}

			t.conn.Write(msg)
		case <-t.stopC:
			return
		}

	}
}
