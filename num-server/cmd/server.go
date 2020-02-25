package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/yurkeen/go-solutions/num-server/deduper"
	"golang.org/x/net/netutil"
)

func handleConnection(conn net.Conn, dd deduper.Deduper) {
	defer conn.Close()
	for {
		rawData, err := bufio.NewReader(conn).ReadSlice('\n')
		if err != nil {
			println("Closing connection:", err.Error())
			return
		}
		println("Received data:", string(rawData))
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		dd.Read(rawData)
	}
}

func run() {
	l, err := net.Listen("tcp", "127.0.0.1:4000")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}
	listener := netutil.LimitListener(l, 5)
	dedup := deduper.New(5)

	for {
		conn, err := listener.Accept()
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn, dedup)
	}
}

func main() {
	run()
}
