// Package httpfromtcp is an HTTP1.1 server that's been built `from scratch`.
// It's a boot.dev project: Learn the HTTP Protocol.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const bufSize = 8

func main() {
	const port = ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error creating listener: port: %s: %s", port, err.Error())
	}
	defer listener.Close()

	log.Printf("Listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: port: %s: %s", port, err.Error())

			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(conn io.ReadCloser) <-chan string {
	buf := make([]byte, bufSize)
	ch := make(chan string)

	go func() {
		defer conn.Close()
		defer close(ch)

		var line strings.Builder

		for {
			n, err := conn.Read(buf)
			if err != nil {
				if line.String() != "" {
					ch <- line.String()
				}
				if errors.Is(err, io.EOF) {
					fmt.Println("connection closed")

					return
				}

				log.Printf("read %v failed: %v", conn, err)

				return
			}

			parts := strings.Split(string(buf[:n]), "\n")
			for _, part := range parts[:len(parts)-1] {
				line.WriteString(part)

				ch <- strings.TrimSpace(line.String())

				line.Reset()
			}

			line.WriteString(parts[len(parts)-1])
		}
	}()

	return ch
}
