package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/bacv/pow-wow/lib/protocol"
	"github.com/bacv/pow-wow/svc"
	"github.com/bacv/pow-wow/svc/wow"
)

func main() {
	log.SetOutput(os.Stdout)
	addr := flag.String("addr", ":8080", "address to receive words of wisdom from")
	flag.Parse()

	wowService := wow.NewWowClientService()

	conn, _ := net.Dial("tcp", *addr)
	transport := svc.NewTransport(conn, wowService.Handle)

	go func() {
		transport.Write(protocol.NewRequestMsg())
	}()

	transport.Spawn()
}
