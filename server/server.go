package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/alsotoes/livelock_simulator/common"
)

var totalMem = 0

func Get(port int, ip string, threadQty int, msgQty int, memMaxPtr int) {

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

	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)

		drop := HandlePackage(threadQueue, remoteaddr, p, memMaxPtr)

		if err != nil {
			log.Printf("**** Some error  %v", err)
			continue
		}

		if drop {
			_, err = conn.WriteToUDP([]byte("-DROP-"), remoteaddr)
		} else {
			_, err = conn.WriteToUDP(p, remoteaddr)
		}

		if err != nil {
			log.Printf("**** Couldn't send response %v", err)
		}
	}
}

func PrepareQueue(threadQty, msgQty int) []*common.Queue {
	threadQueue := make([]*common.Queue, threadQty)

	for i := range threadQueue {
		threadQueue[i] = &common.Queue{NodeStr: make([]*common.Node, msgQty)}
	}

	return threadQueue
}

func HandlePackage(threadQueue []*common.Queue, remoteaddr *net.UDPAddr,
	message []byte, memMaxPtr int) bool {

	drop := false
	msg := string(message[:1024])
	msgArr := strings.Split(msg, "+")

	thread, _ := strconv.Atoi(msgArr[0])
	msgCount, _ := strconv.Atoi(msgArr[1])
	uuid := msgArr[2]

	if totalMem < memMaxPtr {
		totalMem = totalMem + 1
		threadQueue[thread].Push(&common.Node{uuid})
		drop = false
		log.Printf("Read a message from %v [%d,%d] => %s", remoteaddr, thread, msgCount, uuid)
	} else {
		drop = true
		log.Printf("Read a message from %v [%d,%d] => %s", remoteaddr, thread, msgCount, "-DROP-")
	}

	return drop
}

func ForwardingLayer() {
}
