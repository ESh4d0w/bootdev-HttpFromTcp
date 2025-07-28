package request

import (
	"bytes"
	"fmt"
	"io"
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

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine, err := parseRequestLine(req)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	i := bytes.Index(data, []byte(crlf))
	if i == -1 {
		return nil, fmt.Errorf("Couldn't find CRLF")
	}
	requestLineText := string(data[:i])
	requestLine, err := verifyRequestLine(requestLineText)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func verifyRequestLine(reqLine string) (*RequestLine, error) {
	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Request line Didn't contain 3 Parts")
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("Invalid Method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("Maleformed Version: %s", versionParts[:])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("Unrecognized HTTP Version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("Unrecognized HTTP Version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
