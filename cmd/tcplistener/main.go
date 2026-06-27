// -*- mode:go;mode:go-playground -*-
// Copyright © 2026 P, Rich
// License: MIT, see LICENSE for details

// This package is an HTTP/1.1 server built from scratch.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	bufSize = 8
	port    = "42069"
	address = "0:" + port
)

func main() {
	log.SetFlags(0)

	l, err := net.Listen("tcp4", address)
	if err != nil {
		log.Fatalf("error listening on %s: %s", address, err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error on accept: %s", err)

			continue
		}

		fmt.Println("connection accepted")

		for line := range getLinesChannel(conn) {
			log.Println(line)
		}
		fmt.Println("connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	buf := make([]byte, bufSize)

	var s strings.Builder

	go func() {
		defer f.Close()
		defer close(lines)

		for {
			n, err := f.Read(buf)
			if err != nil {
				if n == 0 && errors.Is(err, io.EOF) {
					lines <- s.String()

					return
				}

				log.Printf("error reading: %s\n", err)

				return
			}

			parts := strings.Split(string(buf[:n]), "\n")

			for i, part := range parts {
				if i >= 1 {
					lines <- s.String()

					s.Reset()
				}

				_, err = s.WriteString(part)
				if err != nil {
					log.Printf("error writing to strings.Builder: %s", err)

					break
				}
			}
		}
	}()

	return lines
}
