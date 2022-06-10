package svc

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

var ByteLF = byte(0x0A)

var (
	ErrorSendChannelClosed = errors.New("Send channel is closed")
	ErrorWriteToClosed     = errors.New("Writing to a closed transport")
)

type Request string

type ResponseWriter interface {
	Write(string) error
}

type HandleFunc func(ResponseWriter, Request)

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

	t.Close()
	err := <-errC
	return err
}

func (t *Transport) Close() {
	t.closeOnce.Do(func() {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.closed = true
		close(t.stopC)
	})
}

func (t *Transport) Write(msg string) error {
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
			bytes, err := bufio.NewReader(t.conn).ReadBytes(ByteLF)
			if err != nil {
				errC <- err
				return
			}

			t.handler(t, Request(bytes))
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
