package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	log.Printf("Reading Data from %s:\n ---", inputFilePath)

	for {
		buffer := make([]byte, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("error read: %s\n", err)
			break
		}
		fmt.Printf("read: %s\n", string(buffer[:n]))
	}
}
