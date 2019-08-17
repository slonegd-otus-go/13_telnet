package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	ip := flag.String("ip", "", "ip for listen")
	port := flag.Int("port", 4242, "port for listen")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *ip, *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		connect, err := listener.Accept()
		if err != nil {
			connect.Close()
			continue
		}
		go handle(connect)
	}
}

func handle(connect net.Conn) {
	log.Printf("connect open local: %s, remote %s", connect.LocalAddr(), connect.RemoteAddr())
	scanner := bufio.NewScanner(connect)
	for scanner.Scan() {
		data := scanner.Text()
		log.Println(data)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	log.Printf("connect close local: %s, remote %s", connect.LocalAddr(), connect.RemoteAddr())
}
