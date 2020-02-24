package nums

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strconv"
)

const terminator = `terminate`

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		data, err := bufio.NewReader(conn).ReadSlice('\n')
		if err != nil {
			println("Closing connection:", err)
			return
		}

		dataLen := len(data)
		if dataLen != 9 {
			println("Closing connection: input must be 9 characters, received", dataLen)
			return
		}

		if bytes.Compare(data, []byte(terminator)) == 0 {
			// gracefulStop()
			panic("Terminator string received.\n")
		}

		strData := string(data)

		_, err = strconv.ParseUint(strData, 10, 32)
		// num, err := strconv.ParseUint(strData, 10, 32)
		if err != nil {
			println("Closing connection: input is not a number:", strData)
			return
		}
		println("OK:", strData)

	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:4000")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn)
	}
}
