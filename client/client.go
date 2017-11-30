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

var (
	droppedPackages = 0
	timeoutPackages = 0
	partialTotal    = 0
)

func Call(wg *sync.WaitGroup, counter int, ip string, port int, msgLimit int,
	timeoutPtr int, arrivalRatePtr int, printPtr int) {

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

			if 0 != arrivalRatePtr {
				time.Sleep(time.Millisecond * time.Duration(arrivalRatePtr))
			}

			uuidMsg := common.GenUUID()
			start := time.Now()

			conn.Write([]byte(GenMessage(counter, msgCounter, uuidMsg)))
			response := ProcessResponse(conn)

			t := time.Now()

			partialTotal = partialTotal + 1
			// Dropped packages that could not be push into the queue, so is a loss package
			if "-DROP-" == fmt.Sprintf("%s", response) {
				droppedPackages = droppedPackages + 1
				// Timed out packages, the package is returned but no between the maximun time
			} else if t.Sub(start).Seconds() > float64(timeoutPtr) {
				response = []byte("-TIMEOUT-")
				timeoutPackages = timeoutPackages + 1
			}

			switch printPtr {
			case 1:
				analysedMessages := partialTotal - droppedPackages - timeoutPackages
				csv := fmt.Sprintf("%d,%d", partialTotal, analysedMessages)
				fmt.Println(csv)
			case 2:
				if "-DROP-" != fmt.Sprintf("%s", response) && "-TIMEOUT-" != fmt.Sprintf("%s", response) {
					avgTime := fmt.Sprintf("%f", t.Sub(start).Seconds())
					fmt.Println(avgTime)
				}
			}

			log.Printf("Thread: %d, Msg: %d => Send: %s, Recieved: %s, Elapsed time: %f",
				counter, msgCounter, uuidMsg, response, t.Sub(start).Seconds())

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
