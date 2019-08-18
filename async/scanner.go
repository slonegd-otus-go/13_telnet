package async

import (
	"bufio"
	"io"
	"log"
)

type scanner struct {
	reader io.Reader
	text   chan string
}

func NewScanner(reader io.Reader) scanner {
	return scanner{reader, make(chan string)}
}

func (scanner scanner) Start() {
	go scanner.run()
}

func (s scanner) run() {
	scanner := bufio.NewScanner(s.reader)
	for scanner.Scan() {
		s.text <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func (scanner scanner) Text() <-chan string {
	return scanner.text
}
