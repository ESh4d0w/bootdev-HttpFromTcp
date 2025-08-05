package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/request"
	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/response"
	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerProxy(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody(body)
}
func handlerProxy(w *response.Writer, req *request.Request) {
	count := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + count

	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(0)
	headers.Remove("Content-Length")
	headers.Overwrite("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)

	const maxBufferSize = 1024
	buffer := make([]byte, maxBufferSize)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				log.Printf("Error writing ChunkedBody: %v", err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading resp Body: %v", err)
			break
		}
	}

	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Printf("Error writing ChunkedBodyDone: %v", err)
	}
}
