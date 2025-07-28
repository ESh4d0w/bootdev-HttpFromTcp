package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type parserState int

const (
	parserStateDone        parserState = iota // 0
	parserStateInitialized                    // 1
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		State: parserStateInitialized,
	}

	for req.State != parserStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, cap(buffer)*2)
			_ = copy(newBuffer, buffer)
			buffer = newBuffer
		}
		nRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = parserStateDone
				break
			}
			return nil, err
		}
		readToIndex += nRead

		nParsed, err := req.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[nParsed:])
		readToIndex -= nParsed

	}

	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	i := bytes.Index(data, []byte(crlf))
	if i == -1 {
		return nil, 0, nil
	}

	requestLineText := string(data[:i])
	requestLine, err := verifyRequestLine(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, i + 2, nil
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

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case parserStateInitialized:
		req, read, err := parseRequestLine(data)
		if err != nil {
			return read, fmt.Errorf("Error parsing: %s", err)
		}

		if read == 0 {
			return 0, nil
		}

		r.RequestLine = *req
		r.State = parserStateDone
		return read, nil
	case parserStateDone:
		return 0, fmt.Errorf("Trying to read data in parserStateDone state")
	default:
		return 0, fmt.Errorf("Invalid ParseState: %d", r.State)
	}
}
