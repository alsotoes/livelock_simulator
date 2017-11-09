package client

import (
	"../helper"
	"fmt"
	"log"
	_ "math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Call(wg *sync.WaitGroup, counter int, ip string, port int, msgLimit int) {

	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("udp", addr)
	//conn, err := net.DialTimeout("udp", addr, time.Second)

	defer wg.Done()
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var msg sync.WaitGroup

	for i := 0; i < msgLimit; i++ {
		msg.Add(1)

		go func(msgCounter int) {
			defer msg.Done()

			uuidMsg := uuid.GenUUID()
			start := time.Now()

			conn.Write([]byte(GenMessage(counter, msgCounter, uuidMsg)))
			response := ProcessResponse(conn)

			t := time.Now()

			log.Printf("Counter: %d => Send: %s, Recieved: %s, Elapsed time: %s",
				counter, uuidMsg, response, t.Sub(start))
		}(i)
	}
	msg.Wait()
}

func GenMessage(threadId int, messageId int, uuid string) string {
	message := fmt.Sprintf("%d+%d+%s", threadId, messageId, uuid)
	return message
}

func ProcessResponse(conn net.Conn) []byte {
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

	return buff[:n]
}
