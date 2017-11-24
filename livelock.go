package main

import (
	"flag"
	"log"
	"sync"

	"github.com/alsotoes/livelock_simulator/client"
	"github.com/alsotoes/livelock_simulator/server"
)

func main() {

	impersonationPtr := flag.String("imp", "server", "server or client")
	serverIpPtr := flag.String("ip", "127.0.0.1", "server ip")
	portPtr := flag.Int("port", 3333, "server port to listen for connections")
	threadLimitPtr := flag.Int("threads", 50, "how many threads will be created")
	msgLimitPtr := flag.Int("messages", 50, "how many threads will be created")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	//log.Printf("tail: %s", flag.Args())

	switch *impersonationPtr {
	case "client":
		log.Printf("*** CLIENT CODE")
		var wg sync.WaitGroup

		for i := 0; i < *threadLimitPtr; i++ {
			wg.Add(1)
			go client.Call(&wg, i, *serverIpPtr, *portPtr, *msgLimitPtr)
		}
		wg.Wait()

	case "server":
		log.Printf("*** SERVER CODE")
		server.Get(*portPtr, *serverIpPtr, *threadLimitPtr, *msgLimitPtr)
	default:
		log.Println("*** ERROR: Option unknown")
	}

}
