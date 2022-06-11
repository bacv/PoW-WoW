package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/bacv/pow-wow/lib/protocol"
	"github.com/bacv/pow-wow/svc"
	"github.com/bacv/pow-wow/svc/wow"
)

func main() {
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())
	addr := flag.String("addr", ":8080", "address to receive words of wisdom from")
	flag.Parse()

	wowService := wow.NewWowClientService()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	transport := svc.NewTransport(conn, wowService.Handle)

	go func() {
		transport.Write(protocol.NewRequestMsg())
	}()

	transport.Spawn()
}
