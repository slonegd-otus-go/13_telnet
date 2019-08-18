package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/slonegd-otus-go/13_telnet/async"
)

func main() {
	address := flag.String("adr", "127.0.0.1:4242", "address for connect")
	timeout := flag.Int("timeout", 30, "")
	flag.Parse()

	dialer := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
	connect, err := dialer.DialContext(ctx, "tcp", *address)
	if err != nil {
		log.Fatal(err)
	}

	cancelC := make(chan os.Signal)
	signal.Notify(cancelC, syscall.SIGINT)
	go func() {
		<-cancelC
		log.Println("cancel by interrupt signal")
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		read(ctx, connect)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		write(ctx, connect)
		wg.Done()
	}()

	wg.Wait()
	connect.Close()
}

func read(ctx context.Context, connect net.Conn) {
	scanner := async.NewScanner(connect)
	scanner.Start()
	for {
		select {
		case <-ctx.Done():
			return
		case text := <-scanner.Text():
			log.Println(text)
		}
	}
}

func write(ctx context.Context, connect net.Conn) {
	scanner := async.NewScanner(os.Stdin)
	scanner.Start()
	for {
		select {
		case <-ctx.Done():
			return

		case text := <-scanner.Text():
			text = fmt.Sprintf("%s\n", text)
			_, err := connect.Write([]byte(text))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
