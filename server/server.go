package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alsotoes/livelock_simulator/common"
)

var totalMem = 0

func Get(ip string, port, threadQty, msgQty, memMaxPtr, timeoutPtr int) {

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

		if err != nil {
			log.Printf("**** Some error  %v", err)
			continue
		}

		drop := HandlePackage(threadQueue, remoteaddr, p, memMaxPtr)
		go ForwardingLayer(drop, conn, remoteaddr, p, timeoutPtr)
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
		log.Printf("Read a message from %v [%d,%d] => %s", remoteaddr, thread, msgCount, "-DROP-TIMEOUT-")
	}

	return drop
}

func ForwardingLayer(drop bool, conn *net.UDPConn, remoteaddr *net.UDPAddr,
	message []byte, timeoutPtr int) {

	if drop {
		time.Sleep(time.Second * time.Duration(timeoutPtr))
		message = []byte("-DROP-TIMEOUT-")
	}

	_, err := conn.WriteToUDP(message, remoteaddr)

	if err != nil {
		log.Printf("**** Couldn't send response %v", err)
	}
}
