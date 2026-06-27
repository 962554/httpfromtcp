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

	chunks := readFileInChunks(file)
	for line := range chunksToLines(chunks) {
		log.Println("read:", line)
	}
}

func readFileInChunks(f io.ReadCloser) <-chan string {
	buf := make([]byte, bufSize)
	words := make(chan string)

	go func() {
		defer f.Close()
		defer close(words)

		for {
			n, err := f.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					words <- string(buf[:n])

					return
				}

				log.Printf("error reading file: %s: %s\n", msgFile, err)

				return
			}
			words <- string(buf[:n])
		}
	}()

	return words
}

func chunksToLines(chunks <-chan string) <-chan string {
	lines := make(chan string)

	var s string = ""

	go func() {
		defer close(lines)

		for chunk := range chunks {
			parts := strings.Split(chunk, "\n")

			for i, part := range parts {
				if i >= 1 {
					lines <- s
					s = ""
				}
				s += part
			}
		}
		lines <- s
	}()
	return lines
}
