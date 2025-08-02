package response

import (
	"fmt"
	"io"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header.Set("content-length", fmt.Sprintf("%d", contentLen))
	header.Set("connection", "close")
	header.Set("content-type", "text/plain")
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte(crlf))
	return err
}
