package deduper

import (
	"bytes"
)

const (
	terminator = `terminate`

	// maxBytes == 9 fits uint32; uint64 should be used for longer values.
	maxBytes = 9
)

type Deduper struct {
	seen  map[uint32]bool
	input chan []byte
}

func New(size int) Deduper {
	return Deduper{
		seen:  make(map[uint32]bool),
		input: make(chan []byte, size),
	}
}

func (d *Deduper) Read(line []byte) {
	data := bytes.TrimSpace(line)

	dataLen := len(data)
	if dataLen != maxBytes {
		println("Closing connection: input must be", maxBytes, "characters, received", dataLen)
		return
	}

	if bytes.Compare(data, []byte(terminator)) == 0 {
		gracefulStop()
	}

	if !bytesAreUint(data) {
		println("Closing connection: input is not a number:", string(data))
		return
	}

	println("OK:", string(data))
}

func gracefulStop() {
	panic("Terminator string received.\n")
}

func bytesAreUint(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] < 0x30 || b[i] > 0x39 {
			return false
		}
	}
	return true
}
