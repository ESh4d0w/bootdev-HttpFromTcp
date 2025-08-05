package response

import (
	"fmt"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reason := ""
	switch statusCode {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "InternalServerError"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s%s", statusCode, reason, crlf))
}
