package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alsotoes/livelock_simulator/common"
)

func Call(wg *sync.WaitGroup, counter int, ip string, port int, msgLimit int) {

	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("udp", addr)

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

			uuidMsg := common.GenUUID()
			start := time.Now()

			conn.Write([]byte(GenMessage(counter, msgCounter, uuidMsg)))
			response := ProcessResponse(conn)

			t := time.Now()

			log.Printf("Thread: %d, Msg: %d => Send: %s, Recieved: %s, Elapsed time: %s",
				counter, msgCounter, uuidMsg, response, t.Sub(start))
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
