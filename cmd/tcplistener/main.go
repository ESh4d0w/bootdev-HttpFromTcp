package main

import (
	"fmt"
	"log"
	"net"

	"github.com/esh4d0w/bootdev-HttpFromTcp/internal/request"
)

const inputFilePath = "messages.txt"

func main() {
	tcpListener, err := net.Listen("tcp", "0.0.0.0:42069")
	if err != nil {
		log.Fatalf("Couldn't start tcp listener %v\n", err)
	}
	defer tcpListener.Close()

	for {
		connection, err := tcpListener.Accept()
		if err != nil {
			log.Printf("Connection couldn't be established: %v\n", err)
		}
		log.Printf("Connection opened\n")
		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("Error Getting From Header: %v", err)
		}

		log.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		headers := ""
		for key, value := range req.Headers {
			headers += fmt.Sprintf("- %s: %s\n", key, value)
		}
		log.Printf("Headers:\n%s", headers)

		log.Printf("Connection closed\n")
	}

}
