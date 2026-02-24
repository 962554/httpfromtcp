// Package httpfromtcp is an HTTP1.1 server that's been built `from scratch`.
// It's a boot.dev project: Learn the HTTP Protocol.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	inputFile = "messages.txt"
	bufSize   = 8
)

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("open %s failed: %v", inputFile, err)
	}

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	buf := make([]byte, bufSize)
	ch := make(chan string)

	go func() {
		defer file.Close()
		defer close(ch)

		var line strings.Builder

		for {
			n, err := file.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				log.Printf("read %s failed: %v", inputFile, err)

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
