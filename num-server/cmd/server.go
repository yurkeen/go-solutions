package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/yurkeen/go-solutions/num-server/deduper"
	"golang.org/x/net/netutil"
)

type Config struct {
	maxClients     uint
	port           uint
	host           string
	log            string
	readTimeout    time.Duration
	reportInterval time.Duration
}

func ListenAndServe(conf Config) error {
	stop := make(chan struct{})
	defer close(stop)

	// Initializing network
	nl, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.host, conf.port))
	if err != nil {
		return fmt.Errorf("unable to bind tcp socket: %w", err)
	}
	ll := netutil.LimitListener(nl, int(conf.maxClients))

	ctx, cancel := context.WithCancel(context.Background())

	// Initializing a Deduper
	dd := deduper.New()
	defer dd.Close()

	// Create() will create or truncate the file
	file, err := os.Create(conf.log)
	if err != nil {
		cancel()
		return fmt.Errorf("cannot create log file: %w", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		dd.Log(ctx, file, conf.reportInterval)
		wg.Done()
	}()

	go func() {
		select {
		case <-ctx.Done():
		case <-stop:
		}
		ll.Close()
	}()

	err = serve(ctx, dd, ll, conf.readTimeout, cancel)
	wg.Wait()
	return err
}

// serve accepts connections from net.Listener and passes them
// to deduper.Deduper for futher processing.
//
// Connections have read timeout of time.Duration.
// Parent context can be cancelled using context.CancelFunc.
func serve(ctx context.Context, dd deduper.Deduper, nl net.Listener, rto time.Duration, canc context.CancelFunc) error {

	var (
		wg  sync.WaitGroup
		err error
	)
	defer wg.Wait()

	for {
		conn, e := nl.Accept()
		if e != nil {
			return e
		}
		wg.Add(1)
		go func() {
			err := dd.FromNetwork(conn, rto)
			if err == deduper.ErrTerminationSeq {
				canc()
			}
			conn.Close()
			wg.Done()
		}()
	}
	return err
}

func run(args []string) error {

	// Flags initialization
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	cfg := Config{

		host:           *flags.String("host", "127.0.0.1", "host"),
		port:           *flags.Uint("port", 4000, "port"),
		log:            *flags.String("log", "numbers.log", "log file"),
		maxClients:     *flags.Uint("max", 5, "max clients"),
		reportInterval: *flags.Duration("report", 10*time.Second, "report interval"),
		readTimeout:    *flags.Duration("readTimeout", 10*time.Second, "client read timeout"),
	}
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	return ListenAndServe(cfg)
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
