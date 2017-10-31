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
	_ "time"
)

const (
	Message       = "pong"
)

func Get(port int, ip string) {

    addr := net.UDPAddr{
        Port: port,
        IP: net.ParseIP(ip),
    }

    p := make([]byte, 1024)
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %d", port)

	for {
        _,remoteaddr,err := conn.ReadFromUDP(p)
        log.Printf("Read a message from %v %s", remoteaddr, p)

        if err !=  nil {
            log.Printf("**** Some error  %v", err)
            continue
        }
        go handler(conn, remoteaddr, p)
	}

}

func handler(conn *net.UDPConn, addr *net.UDPAddr, message []byte) {
    _,err := conn.WriteToUDP([]byte(Message), addr)
    //_,err := conn.WriteToUDP(message, addr)

    if err != nil {
        log.Printf("**** Couldn't send response %v", err)
    }
}
