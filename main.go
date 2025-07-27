package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	log.Printf("Reading Data from %s:\n ---", inputFilePath)

	currentLine := ""
	for {
		buffer := make([]byte, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
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
		// if len(parts) > 1 {
		// 	fmt.Printf("read:%s%s\n", currentLine, parts[0])
		// 	currentLine = ""
		// }
		// Commented is my solution the one below is the from boot.dev
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
}
