package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/bacv/pow-wow/svc"
	"github.com/bacv/pow-wow/svc/wow"
)

func main() {
	log.SetOutput(os.Stdout)
	port := flag.Int("port", 8080, "port to run words of wisdom server on")
	flag.Parse()

	wowSource := wow.NewWisdomSource()
	wowService := wow.NewWowServerService(wowSource)
	tcpServer := svc.NewTcpServer(wowService.Handle)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := tcpServer.Serve(*port)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// TODO: gracefull shutdown
	wg.Wait()
}
