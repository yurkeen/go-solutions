package deduper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	terminator = `terminate`

	// maxBytes == 9 fits uint32; uint64 should be used for longer values.
	maxBytes = 9
)

type Deduper struct {
	file  os.File
	seen  map[uint32]bool
	input chan []byte
}

func New(log *os.File, size int) Deduper {

	d := Deduper{
		file:  *os.file,
		seen:  make(map[uint32]bool),
		input: make(chan []byte, size),
	}
	return d
}

func (d *Deduper) Ingest(line []byte) error {
	data := bytes.TrimSpace(line)

	dataLen := len(data)
	if dataLen != maxBytes {
		return fmt.Errorf("expected %d characters, received %d", maxBytes, dataLen)
	}

	if bytes.Compare(data, []byte(terminator)) == 0 {
		close(d.input) // This may panic, use ctx
		return nil
	}

	if !bytesAreUint(data) {
		return fmt.Errorf("input is not a number: %s", string(data))
	}
	d.input <- data
	return nil
}

func (d *Deduper) Log() {

	for data := range d.input {
		number := binary.BigEndian.Uint32(data)
		if _, ok := d.seen[number]; !ok {
			d.seen[number] = true
			d.file.Write(append(data, '\n'))
		}
	}
	return
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
