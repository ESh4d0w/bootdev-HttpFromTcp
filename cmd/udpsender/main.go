package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalf("Error Resolving UDP Addr: %v\n", err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Error Dial TCP: %v\n", err)
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input %v\n", err)
		}
		_, err = udpConn.Write([]byte(message))
		if err != nil {
			log.Fatalf("Error sending message: %v\n", err)
		}
		fmt.Printf("Message sent: %s\n", message)
	}

}
