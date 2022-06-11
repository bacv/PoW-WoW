package svc

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	handleFunc HandleFunc
}

func NewTcpServer(handleFunc HandleFunc) *Server {
	return &Server{
		handleFunc: handleFunc,
	}
}

func (s *Server) Serve(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			tErr := s.HandleTcp(conn)
			if tErr != nil {
				log.Print(err)
			}
		}()
	}
}

func (s *Server) HandleTcp(conn net.Conn) error {
	transport := NewTransport(conn, s.handleFunc)

	return transport.Spawn()
}
