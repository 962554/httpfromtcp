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
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var methods = map[string]any{
	"GET":  struct{}{},
	"POST": struct{}{},
	"HEAD": struct{}{},
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	mesg := string(b)
	parts := strings.Split(mesg, "\r\n")
	reqLine := parts[0]
	log.Println("parsing:", reqLine)

	reqParts := strings.Split(reqLine, " ")
	// ensure request-line has 3 parts
	if len(reqParts) != 3 {
		return nil, fmt.Errorf("malformed request-line: got %d parts, expecting: 3 parts", len(reqParts))
	}

	method, target, httpVersion := reqParts[0], reqParts[1], reqParts[2]
	if _, ok := methods[method]; !ok {
		return nil, fmt.Errorf("unknown method: %s", method)
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("target: %s does not start with /", target)
	}

	versionParts := strings.Split(httpVersion, "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed version: there should only be one /")
	}
	prefix, version := versionParts[0], versionParts[1]
	if prefix != "HTTP" {
		return nil, fmt.Errorf("malformed version: got: %s, expected: HTTP", prefix)
	}

	if version != "1.1" {
		return nil, fmt.Errorf("malformed version: got: %s, expected: 1.1", version)
	}

	return &Request{RequestLine{HttpVersion: version, RequestTarget: target, Method: method}}, nil
}
