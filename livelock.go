package main

import (
	"flag"
	"io/ioutil"
	"log"
	"sync"

	"github.com/alsotoes/livelock_simulator/client"
	"github.com/alsotoes/livelock_simulator/server"
)

func LogStdOut(logApp bool) {

	if logApp {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	} else {
		log.SetOutput(ioutil.Discard)
	}

}

func main() {

	impersonationPtr := flag.String("imp", "server", "server or client")
	serverIpPtr := flag.String("ip", "127.0.0.1", "server ip")
	portPtr := flag.Int("port", 3333, "server port to listen for connections")
	threadLimitPtr := flag.Int("threads", 50, "how many threads will be created")
	msgLimitPtr := flag.Int("messages", 50, "how many threads will be created")
	memMaxPtr := flag.Int("memory", 1500, "maximun global memory to store messages")
	timeoutPtr := flag.Int("timeout", 1, "timeout in seconds to drop packages")
	debugPtr := flag.Bool("debug", true, "log to stdout")
	flag.Parse()

	LogStdOut(*debugPtr)

	if (*threadLimitPtr)*(*msgLimitPtr) < (*memMaxPtr) {
		*memMaxPtr = (*threadLimitPtr) * (*msgLimitPtr)
	}

	switch *impersonationPtr {
	case "client":
		log.Printf("*** CLIENT CODE")
		var wg sync.WaitGroup

		for i := 0; i < *threadLimitPtr; i++ {
			wg.Add(1)
			go client.Call(&wg, i, *serverIpPtr, *portPtr, *msgLimitPtr, *timeoutPtr)
		}
		wg.Wait()

	case "server":
		log.Printf("*** SERVER CODE")
		server.Get(*serverIpPtr, *portPtr, *threadLimitPtr, *msgLimitPtr,
			*memMaxPtr, *timeoutPtr)
	default:
		log.Println("*** ERROR: Option unknown")
	}

}
