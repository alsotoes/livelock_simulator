package client

import (
	_ "fmt"
	"log"
	_ "math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	StopCharacter = "\r\n\r\n"
)

func Call(wg *sync.WaitGroup, counter int, ip string, port int, uuid string) {

	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp4", addr)

	defer wg.Done()
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()

	conn.Write([]byte(uuid))
	conn.Write([]byte(StopCharacter))

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

	t := time.Now()

	log.Printf("Counter: %d => Send: %s, %s, %s", counter, uuid, buff[:n], t.Sub(start))
}
