package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yurkeen/go-solutions/num-server/deduper"
	"golang.org/x/net/netutil"
)

var ErrNetClosing = errors.New("use of closed network connection")

type Server struct {
	listener net.Listener
	stop     chan struct{}
}

func StartServer(ctx context.Context, host string, port int, clients int) (*Server, error) {
	// Initializing network
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, fmt.Errorf("unable to bind tcp socket: %w", err)
	}
	listener := netutil.LimitListener(l, clients)
	s := &Server{
		listener: listener,
		stop:     make(chan struct{}),
	}
	return s, nil
}

func (s *Server) serve(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
		}
		conn, err := s.listener.Accept()
		if err != nil {
			if !errors.Is(err, ErrNetClosing) {
				fmt.Fprintf(os.Stderr, "error accepting TCP connection: %s", err.Error())
			}
			break
		}
		connCh <- conn
	}
	return
}

func run(args []string) error {

	// Flags initialization
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		host           = *flags.String("host", "127.0.0.1", "host")
		port           = *flags.Int("port", 4000, "port")
		log            = *flags.String("log", "numbers.log", "log file")
		maxClients     = *flags.Int("max", 5, "max clients")
		reportInterval = *flags.Duration("report", 10*time.Second, "report interval")
		readTimeout    = *flags.Duration("readTimeout", 10*time.Second, "client read timeout")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// Create() will create or truncate the file
	file, err := os.Create(log)
	if err != nil {
		return fmt.Errorf("cannot create log file: %w", err)
	}
	defer file.Close()

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize a deduper
	dd := deduper.New()
	defer dd.Close()

	// Start writing log and reporting to console
	go dd.Log(ctx, file, reportInterval)

	var (
		wg     sync.WaitGroup
		runErr error
	)
	// Handling connections until context cancelled
	connCh := make(chan net.Conn, 1)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				break
			default:
			}
			println("new connection waiting")
			conn, err := listener.Accept()
			if err != nil {
				if !errors.Is(err, ErrNetClosing) {
					runErr = fmt.Errorf("error accepting TCP connection: %w", err)
				}
				println("aborted new connection")
				break
			}
			println("connection accepted")
			connCh <- conn
		}
		println("connection dispatcher returning")
		return
	}(ctx)

	for conn := range connCh {
		wg.Add(1)
		go func() {
			conn.SetReadDeadline(time.Now().Add(readTimeout))

			err := dd.FromNetwork(conn, readTimeout)
			if err == deduper.ErrTerminationSeq {
				cancel()
				close(connCh)
				println("conn channel closed")
			}
			conn.Close()
			wg.Done()
			return
		}()
	}
	println("waiting for all goroutines to complete")
	wg.Wait()
	return runErr
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
