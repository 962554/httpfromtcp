// -*- mode:go;mode:go-playground -*-
// Copyright © 2026 P, Rich
// License: MIT, see LICENSE for details

// This package handles requests.

package request

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	Line Line
}

type Line struct {
	HTTPVersion string
	Target      string
	Method      string
}

func FromReader(reader io.Reader) (*Request, error) {
	mesg, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading all: %w", err)
	}

	reqLine, err := LineParse(string(mesg))
	if err != nil {
		return nil, err
	}

	return &Request{Line: *reqLine}, nil
}

func LineParse(mesg string) (*Line, error) {
	methods := map[string]any{
		"GET":  struct{}{},
		"POST": struct{}{},
		"HEAD": struct{}{},
	}

	const (
		lineParts    = 3
		versionParts = 2
	)

	parts := strings.Split(mesg, "\r\n")
	reqLine := parts[0]
	log.Println("parsing:", reqLine)

	reqParts := strings.Split(reqLine, " ")
	// ensure request-line has 3 parts
	if len(reqParts) != lineParts {
		return nil, fmt.Errorf("malformed request-line: got %d parts, expecting: %d parts", len(reqParts), lineParts)
	}

	method, target, httpVersion := reqParts[0], reqParts[1], reqParts[2]
	if _, ok := methods[method]; !ok {
		return nil, fmt.Errorf("unknown method: %s", method)
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("target: %s does not start with /", target)
	}

	verParts := strings.Split(httpVersion, "/")
	if len(verParts) != versionParts {
		return nil, fmt.Errorf("malformed version: got %d parts, expecting: %d parts", len(verParts), versionParts)
	}

	prefix, version := verParts[0], verParts[1]
	if prefix != "HTTP" {
		return nil, fmt.Errorf("malformed version: got: %s, expected: HTTP", prefix)
	}

	if version != "1.1" {
		return nil, fmt.Errorf("malformed version: got: %s, expected: 1.1", version)
	}

	return &Line{HTTPVersion: version, Target: target, Method: method}, nil
}
