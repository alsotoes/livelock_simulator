package main

import (
	"./client"
	"./helper"
	"./server"
	"flag"
	"log"
	"sync"
)

var (
	ip    string = "127.0.0.1"
	port  int    = 3333
	limit int    = 250
)

func main() {

	impersonationPtr := flag.String("imp", "server", "Server or client")
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
			go client.Call(&wg, i, ip, port, uuid.GenUUID())
		}
		wg.Wait()

	case "server":
		log.Printf("*** SERVER CODE")
		server.Get(port, ip)
	default:
		log.Println("*** ERROR: Option unknown")
	}

}
