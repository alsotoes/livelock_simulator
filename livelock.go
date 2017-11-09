package main

import (
	"./client"
	"./server"
	"flag"
	"log"
	"sync"
)

var (
	ip          string = "127.0.0.1"
	port        int    = 3333
	limit       int    = 50
    msgLimit    int    = 50
)

func main() {

	impersonationPtr := flag.String("imp", "server", "server or client")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	//log.Printf("impersonation: %s", *impersonationPtr)
	//log.Printf("tail: %s", flag.Args())

	switch *impersonationPtr {
	case "client":
		log.Printf("*** CLIENT CODE")
		var wg sync.WaitGroup

		for i := 0; i < limit; i++ {
			wg.Add(1)
			go client.Call(&wg, i, ip, port, msgLimit)
		}
		wg.Wait()

	case "server":
		log.Printf("*** SERVER CODE")
		server.Get(port, ip)
	default:
		log.Println("*** ERROR: Option unknown")
	}

}
