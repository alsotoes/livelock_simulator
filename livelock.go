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
	msgLimitPtr := flag.Int("messages", 50, "how many messages will be sended")
	memMaxPtr := flag.Int("memory", 1500, "maximun global memory to store messages")
	timeoutPtr := flag.Int("timeout", 10, "timeout in seconds to drop packages")
	arrivalRatePtr := flag.Int("rate", 0, "arrival rate in miliseconds (default 0)")
	printPtr := flag.Int("print", 0, "print csv to stdout 0-noprint 1-packetrate 2-time 3-dropped 4-timeout (default 0)")
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
			go client.Call(&wg, i, *serverIpPtr, *portPtr, *msgLimitPtr, *timeoutPtr,
				*arrivalRatePtr, *printPtr)
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
