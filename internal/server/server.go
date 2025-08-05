package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/request"
	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error Accepting con: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	resWriter := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		resWriter.WriteStatusLine(response.StatusBadRequest)
		body := fmt.Appendf([]byte{}, "Error Reading request: %v", err)
		resWriter.WriteHeaders(response.GetDefaultHeaders(len(body)))
		resWriter.WriteBody(body)
		return
	}
	s.handler(resWriter, req)
}
