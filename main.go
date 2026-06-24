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
	defer file.Close()

	buf := make([]byte, bufSize)

	var s strings.Builder

	for {
		n, err := file.Read(buf)
		if err != nil {
			if n == 0 && errors.Is(err, io.EOF) {
				log.Printf("read: %s\n", s.String())

				break
			}

			log.Printf("error reading file: %s: %s\n", msgFile, err)

			break
		}

		parts := strings.Split(string(buf[:n]), "\n")

		for i, part := range parts {
			if i >= 1 {
				log.Printf("read: %s\n", s.String())
				s.Reset()
			}

			_, err = s.WriteString(part)
			if err != nil {
				log.Printf("error writing to strings.Builder: %s", err)

				break
			}
		}
	}
}
