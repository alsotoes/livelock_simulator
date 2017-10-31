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

func Call(wg *sync.WaitGroup, counter int, ip string, port int, uuid string) {

	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("udp", addr)

	defer wg.Done()
	defer conn.Close()

	if err != nil {
		log.Fatalln(err)
	}

	start := time.Now()

	conn.Write([]byte(uuid))

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)

	t := time.Now()

	log.Printf("Counter: %d => Send: %s, Recieved: %s, Elapsed time: %s",
        counter, uuid, buff[:n], t.Sub(start))
}
