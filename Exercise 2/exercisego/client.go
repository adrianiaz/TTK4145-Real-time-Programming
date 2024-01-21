package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func client() {
	connection, err := net.Dial("tcp4", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	//loop for reading/writing to the server
	for {
		textIn = bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
	}
}
