package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/yurkeen/go-solutions/num-server/deduper"
	"golang.org/x/net/netutil"
)

func handleConnections(nl net.Listener, dd deduper.Deduper) {
	for {
		conn, err := nl.Accept()
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		if err != nil {
			println("error accepting TCP connection:", err.Error())
			return
		}

		go func() {
			defer conn.Close()
			for {
				rawData, err := bufio.NewReader(conn).ReadSlice('\n')
				if err != nil {
					println("Closing connection:", err.Error())
					return
				}
				println("Received data:", string(rawData))
				conn.SetReadDeadline(time.Now().Add(30 * time.Second))

				if err := dd.Ingest(rawData); err != nil {
					println("Closing connection:", err.Error())
					return
				}
			}
		}()
	}
}

func handleLog(dd deduper.Deduper) {

}

func run() {
	l, err := net.Listen("tcp", "127.0.0.1:4000")
	if err != nil {
		log.Fatal("tcp server listener error:", err.Error())
	}
	listener := netutil.LimitListener(l, 5)
	defer listener.Close()

	// Create() will create or truncate the file
	file, err := os.Create("numbers.log")
	if err != nil {
		log.Fatal("cannot create log file:", err.Error())
	}
	defer file.Close()

	dedup := deduper.New(file, 5)

	go handleConnections(listener, dedup)
	dedup.Log()
}

func main() {
	run()
}
