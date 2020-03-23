package deduper

import (
	"bufio"
	"bytes"
	"context"
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

// payload is used to transmit the number alongside
// its []byte representation to the Log() goroutine.
type payload struct {
	number uint32
	bytes  []byte
}

// Deduper receives slices of bytes as messsages
// and after processing writes those to
// a deduplicated stream.
type Deduper struct {
	// input receives vetted slices for
	// further processing
	input chan payload

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
		input: make(chan payload),
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
		fmt.Println("Received data:", data)
		if bytes.Compare(data, []byte(terminator)) == 0 {
			err = ErrTerminationSeq
			break
		}

		if len(data) != expectedBytes {
			err = ErrUnexpectedInput
			break
		}
		var num uint32
		if err = bytesToUint32(data, &num); err != nil {
			break
		}

		// Both number and byte data sent down the input channel
		// to avoid reverse conversion from uint32 to []byte.
		d.input <- payload{num, data}
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
	defer func() {
		println("Flushing file")
		bw.Flush()
	}()

	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()
	for {
		select {
		case data := <-d.input:
			fmt.Printf("Number: %d\n", data.number)
			if d.lookup(data.number) {
				dups.inc()
				break
			}
			d.record(data.number)
			uniq.inc()

			println("Writing:", data.bytes)
			_, d.err = bw.Write(append(data.bytes, '\n'))
			if d.err != nil {
				println("Error writing file:", d.err)
				// No need to flush, can return immediately
				return
			}
		case <-ticker.C:
			fmt.Printf(
				"Received %d unique numbers, %d duplicates. Unique total: %d\n",
				uniq, dups, d.count(),
			)
			uniq.reset()
			dups.reset()
		case <-ctx.Done():
			return
		}
	}
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

func bytesToUint32(b []byte, n *uint32) error {

	var (
		pos uint32 = 1
		num uint32
	)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] < '0' || b[i] > '9' {
			return errors.New("byts are not numbers")
		}
		num += pos * (uint32(b[i]) - '0')
		pos *= 10 // base 10
	}
	*n = num
	return nil
}
