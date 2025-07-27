package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesChannel := getLinesChannel(connection)
		for line := range linesChannel {
			log.Printf("%s\n", line)
		}
		log.Printf("Connection closed\n")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				log.Printf("error read: %s\n", err)
				break
			}
			currentStr := string(buffer[:n])
			parts := strings.Split(currentStr, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return lines
}
