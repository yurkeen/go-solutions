package deduper

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	// 9 fits uint32, otherwise uint64 should be used
	// for longer values.
	expectedBytes = 9
	terminator    = `terminate`
)

var (
	ErrReadTimeout     = errors.New("read timeout")
	ErrUnexpectedInput = errors.New("unexpected input")
	ErrTerminationSeq  = errors.New("termination sequence received")
)

// Deduper receives slices of bytes as messsages
// and after processing writes those to
// a deduplicated stream.
type Deduper struct {
	// input receives vetted slices for
	// further processing
	input chan []byte

	// seen is a lookup map to deduplicate
	// values
	seen map[uint32]struct{}

	// error to be returned when processing
	// of data is completed
	err error
}

// New returns an initialized Deduper.
func New() Deduper {
	return Deduper{
		input: make(chan []byte, 1),
		seen:  make(map[uint32]struct{}),
	}
}

// FromNetwork validates input bytes, collected from net.Conn and writes
// validated stream of []byte into the input channel. An error will be
// returned if validation does not pass or io.EOF.
//
// Before being written to the input channel, data is validated to be
// exactly expectedBytes length and to represent an uint32 number.
//
// Any further processing of the stream will be immediately stopped upon
// receiving a termination string and the ErrTerminationSeq will be returned.
//
// An io.EOF will be returned as error once no more bytes left to read from the
// source io.Reader.
func (d Deduper) FromNetwork(nc net.Conn, readTimeout time.Duration) error {
	var (
		err      error
		byteData []byte
	)
	for {
		byteData, err = bufio.NewReader(nc).ReadSlice('\n')
		if err != nil {
			break
		}
		nc.SetReadDeadline(time.Now().Add(readTimeout))

		data := bytes.TrimSpace(byteData)
		if bytes.Compare(data, []byte(terminator)) == 0 {
			err = ErrTerminationSeq
			break
		}

		if len(data) != expectedBytes || !bytesAreUint(data) {
			err = ErrUnexpectedInput
			break
		}

		d.input <- data
	}
	return err
}

// Log writes deduplicated stream of byte slices from the input
// and reports stats to console.
//
// Newline '\n' is used as a separator for writing to the specified
// io.Writer. Stats are written every time.Duration.
// Any error while writing the data will be recorded in d.err and
// the method will return.
func (d *Deduper) Log(ctx context.Context, w io.Writer, reportInterval time.Duration) {

	// Counters for stats reporting
	var uniq, dups counter

	// Buffered writes can help performance during the Write() operations
	// by buffering 4096 bytes before writing.
	bw := bufio.NewWriter(w)

	for {
		select {
		case data := <-d.input:

			println("Data:", data)
			number := binary.BigEndian.Uint32(data)
			println("Number:", number)
			if d.lookup(number) {
				dups.inc()
				continue
			}

			d.record(number)
			uniq.inc()

			_, d.err = bw.Write(append(data, '\n'))
			if d.err != nil {
				// No need to flush, can return
				// immediately
				return
			}
		case <-time.Tick(reportInterval):
			fmt.Printf(
				"Received %d unique numbers, %d duplicates. Unique total: %d\n",
				uniq, dups, d.count(),
			)
			uniq.reset()
			dups.reset()
		case <-ctx.Done():
			break
		}
	}
	d.err = bw.Flush()
	return
}

// Close closes the deduper.
//
// It returns any error encountered during the close operation.
func (d *Deduper) Close() error {
	println("Close invoked")
	close(d.input)
	return d.err
}

// seen returns whether a number was seen and recorded
// by Deduper.
func (d Deduper) lookup(n uint32) bool {
	if _, ok := d.seen[n]; ok {
		return true
	}
	return false
}

// record checks in the given number for further lookup.
func (d *Deduper) record(n uint32) {
	d.seen[n] = struct{}{}
}

func (d Deduper) count() int {
	return len(d.seen)
}

// Counter represents a single forward-only counter.
type counter uint

// inc increments the counter by 1.
func (c *counter) inc() { *c++ }

// reset resets the counter back to 0.
func (c *counter) reset() { *c = 0 }

func bytesAreUint(b []byte) bool {
	for i := 0; i < len(b); i++ {
		if b[i] < 0x30 || b[i] > 0x39 {
			return false
		}
	}
	return true
}
