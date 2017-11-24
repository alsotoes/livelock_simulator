package server

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/alsotoes/livelock_simulator/common"
)

func Get(port int, ip string, threadQty int, msgQty int) {

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	p := make([]byte, 1024)
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen => %s:%d", ip, port)

	threadQueue := PrepareQueue(threadQty, msgQty)

	threadQueue[0].Push(&common.Node{"HOLA!"})
	log.Printf("== VALUE ==> %s", threadQueue[0].Pop().Value)

	// INICIO: test de codigo
	/*
		go func() {
			queue[50] <- "asddasd"
			queue[50] <- "xxxxxxx"
			queue[50] <- "yyyyyyy"
			log.Printf("=> %d", len(queue[50]))
		}()

		x := <-queue[50]
		log.Printf("%s", x)

		log.Printf("Begin listen port: %d", port)
	*/
	// FIN: test de codigo

	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		//HandlePackage(queue, p)
		//log.Printf("Read a message from %v %s", remoteaddr, p)

		if err != nil {
			log.Printf("**** Some error  %v", err)
			continue
		}

		_, err = conn.WriteToUDP(p, remoteaddr)
		if err != nil {
			log.Printf("**** Couldn't send response %v", err)
		}
	}

}

func PrepareQueue(threadQty, msgQty int) []*common.Queue {
	threadQueue := make([]*common.Queue, threadQty)

	for i := range threadQueue {
		threadQueue[i] = &common.Queue{Nodex: make([]*common.Node, msgQty)}
	}

	return threadQueue
}

func HandlePackage(queue [100]chan string, message []byte) {
	msg := string(message[:1024])
	msgArr := strings.Split(msg, "+")

	//thread := msgArr[0]
	//msgCount := msgArr[1]
	uuid := msgArr[2]

	log.Printf("Read a message from %s", uuid)
	/*
	   go func() {
	       queue[50] <- "asddasd"
	       queue[50] <- "xxxxxxx"
	       queue[50] <- "yyyyyyy"
	       log.Printf("=> %d", len(queue[50]))
	   }()
	*/
}
