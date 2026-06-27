// -*- mode:go;mode:go-playground -*-
// Copyright © 2026 P, Rich
// License: MIT, see LICENSE for details

// This package is an HTTP/1.1 server built from scratch.
package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

const (
	msgFile = "messages.txt"
	bufSize = 8
)

func main() {
	log.SetFlags(0)

	file, err := os.Open(msgFile)
	if err != nil {
		log.Fatalf("error opening file: %s for reading: %s", msgFile, err)
	}

	for line := range getLinesChannel(file) {
		log.Println("read:", line)
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
