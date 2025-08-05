package response

import (
	"fmt"
	"io"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/headers"
)

const crlf = "\r\n"

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	state  writerState
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		state:  writerStateStatusLine,
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return fmt.Errorf("cannot write StatusLine in state :%d", w.state)
	}
	defer func() { w.state = writerStateHeaders }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return fmt.Errorf("cannot write Headers in state :%d", w.state)
	}
	defer func() { w.state = writerStateBody }()
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s%s", key, value, crlf)
		_, err := w.writer.Write([]byte(header))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(crlf))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write Body in state :%d", w.state)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write Body in state :%d", w.state)
	}

	chunkSize := len(p)
	total := 0
	n, err := fmt.Fprintf(w.writer, "%x%s", chunkSize, crlf)
	if err != nil {
		return total, err
	}
	total += n

	n, err = w.writer.Write(p)
	if err != nil {
		return total, err
	}
	total += n

	n, err = fmt.Fprintf(w.writer, "%s", crlf)
	if err != nil {
		return total, err
	}
	total += n

	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("cannot write Body in state :%d", w.state)
	}
	n, err := fmt.Fprintf(w.writer, "0%s%s", crlf, crlf)
	if err != nil {
		return n, err
	}
	return n, nil

}
