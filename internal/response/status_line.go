package response

import (
	"fmt"
	"io"
)

const crlf = "\r\n"

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reason := ""
	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "InternalServerError"
	}

	line := fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, reason, crlf)

	_, err := w.Write([]byte(line))
	if err != nil {
		return err
	}
	return nil
}
