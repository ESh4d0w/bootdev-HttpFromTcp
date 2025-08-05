package response

import (
	"fmt"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header.Set("content-length", fmt.Sprintf("%d", contentLen))
	header.Set("connection", "close")
	header.Set("content-type", "text/plain")
	return header
}
