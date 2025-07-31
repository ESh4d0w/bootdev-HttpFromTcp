package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/headers"
)

type requestState int

const (
	requestStateDone           requestState = iota // 0
	requestStateInitialized                        // 1
	requestStateParsingHeaders                     // 2
	requestStateParsingBody                        // 3
)

type Request struct {
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	BodyLengthRead int
	State          requestState
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
		State:   requestStateInitialized,
		Headers: headers.NewHeaders(),
		Body:    make([]byte, 0),
	}

	for req.State != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, cap(buffer)*2)
			_ = copy(newBuffer, buffer)
			buffer = newBuffer
		}
		nRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.State, nRead)
				}
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
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateParsingBody:
		contentLenString, ok := r.Headers.Get("Content-Length")
		if !ok {
			// Ignore if no content-length
			r.State = requestStateDone
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(contentLenString)
		if err != nil {
			return 0, fmt.Errorf("Malformed Content-Length: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.BodyLengthRead += len(data)
		if r.BodyLengthRead > contentLen {
			return 0, fmt.Errorf("Body longer than Content-Length")
		}
		if r.BodyLengthRead == contentLen {
			r.State = requestStateDone
		}
		return len(data), nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Error parsing Header: %s", err)
		}
		if done {
			r.State = requestStateParsingBody
		}
		return n, nil
	case requestStateInitialized:
		req, n, err := parseRequestLine(data)
		if err != nil {
			return n, fmt.Errorf("Error parsing Request Line: %s", err)
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *req
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("Trying to read data in requestStateDone state")
	default:
		return 0, fmt.Errorf("Invalid ParseState: %d", r.State)
	}
}
