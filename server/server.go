package server

import (
	_ "bufio"
	_ "io"
	"log"
	_ "math/rand"
	"net"
	"os"
	_ "strconv"
	_ "strings"
	_ "syscall"
	_ "time"
)

/*
func setNoReuseAddress(conn net.PacketConn) {
	file, _ := conn.(*net.UDPConn).File()
	fd := file.Fd()
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0)
}
*/

func Get(port int, ip string, threadQty int, msgQty int) {

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	p := make([]byte, 1024)
	conn, err := net.ListenUDP("udp", &addr)
	//setNoReuseAddress(conn)
	//conn.SetTimeout(500 * 1000 * 1000)

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen => %s:%d", ip, port)

	queue := PrepareQueue(msgQty)

	// INICIO: test de codigo
	go func() {
		queue[50] <- "asddasd"
		queue[50] <- "xxxxxxx"
		queue[50] <- "yyyyyyy"
		log.Printf("=> %d", len(queue[50]))
	}()

	x := <-queue[50]
	log.Printf("%s", x)

	log.Printf("Begin listen port: %d", port)
	// FIN: test de codigo

	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		//HandlePackage(queue, p)
		log.Printf("Read a message from %v %s", remoteaddr, p)

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

func PrepareQueue(msgQty int) [100]chan string {
	var queue [100]chan string

	for i := range queue {
		queue[i] = make(chan string, msgQty)
	}

	return queue
}

func HandlePackage(queue chan string, message []byte) {
}
