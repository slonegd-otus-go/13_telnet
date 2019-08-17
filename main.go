package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func unused(interface{}) {}

func main() {
	dialer := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	connect, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:4242")
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

type scanner struct {
	reader io.Reader
	text   chan string
}

func newScanner(reader io.Reader) scanner {
	return scanner{reader, make(chan string)}
}

func (s scanner) Run() {
	scanner := bufio.NewScanner(s.reader)
	for scanner.Scan() {
		data := scanner.Text()
		s.text <- data
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func read(ctx context.Context, connect net.Conn) {
	scanner := newScanner(connect)
	go scanner.Run()
	for {
		select {
		case <-ctx.Done():
			return
		case text := <-scanner.text:
			log.Println(text)
		}
	}
}

func write(ctx context.Context, connect net.Conn) {
	scanner := newScanner(os.Stdin)
	go scanner.Run()
	for {
		select {
		case <-ctx.Done():
			return

		case text := <-scanner.text:
			text = fmt.Sprintf("%s\n", text)
			_, err := connect.Write([]byte(text))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
