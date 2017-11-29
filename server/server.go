package server

import (
	"log"
	"math/rand"
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

		thread, msgCount := HandlePackage(threadQueue, remoteaddr, p, memMaxPtr)

		go func(conn *net.UDPConn, remoteaddr *net.UDPAddr, threadQueue []*common.Queue,
			thread int, msgCount int, p []byte, memMaxPtr int) {

			message := ForwardingLayer(conn, remoteaddr, threadQueue, thread, msgCount, timeoutPtr)
			ExecutePetition(1500)

			log.Printf("Read a message from %v [%d,%d] => %s", remoteaddr, thread, msgCount, message)
			_, err = conn.WriteToUDP(message, remoteaddr)

			if err != nil {
				log.Printf("**** Couldn't send response %v", err)
			}

		}(conn, remoteaddr, threadQueue, thread, msgCount, p, memMaxPtr)
	}
}

func BlackBox(conn *net.UDPConn, remoteaddr *net.UDPAddr, threadQueue []*common.Queue,
	p []byte, memMaxPtr int) {
}

// from this line and till the end the code need to be put in other
// file inside server package

func PrepareQueue(threadQty, msgQty int) []*common.Queue {
	threadQueue := make([]*common.Queue, threadQty)

	for i := range threadQueue {
		threadQueue[i] = &common.Queue{NodeStr: make([]*common.Node, msgQty)}
	}

	return threadQueue
}

func HandlePackage(threadQueue []*common.Queue, remoteaddr *net.UDPAddr,
	message []byte, memMaxPtr int) (int, int) {

	msg := string(message[:1024])
	msgArr := strings.Split(msg, "+")

	thread, _ := strconv.Atoi(msgArr[0])
	msgCount, _ := strconv.Atoi(msgArr[1])
	uuid := msgArr[2]

	if totalMem < memMaxPtr {
		totalMem = totalMem + 1
		threadQueue[thread].Push(&common.Node{uuid})
	}

	return thread, msgCount
}

func ForwardingLayer(conn *net.UDPConn, remoteaddr *net.UDPAddr,
	threadQueue []*common.Queue, thread int, msgCount int, timeoutPtr int) []byte {

	message := []byte("")

	if indexCheck := threadQueue[thread]; indexCheck != nil {
		if rawValue := indexCheck.Pop(); rawValue != nil {
			message = []byte(rawValue.Value)
		} else {
			message = []byte("-DROP-")
		}
	}

	totalMem = totalMem - 1
	/*
		////
		rand.Seed(time.Now().Unix())
		////time.Sleep(time.Millisecond * time.Duration(rand.Intn(fixRandom)))
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		//ExecutePetition(1500)

		_, err := conn.WriteToUDP(message, remoteaddr)
		log.Printf("Read a message from %v [%d,%d] => %s", remoteaddr, thread, msgCount, message)

		if err != nil {
			log.Printf("**** Couldn't send response %v", err)
		}
	*/

	return message
}

func ExecutePetition(fixRandom int) {
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(fixRandom)))
}
