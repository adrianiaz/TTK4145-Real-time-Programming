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
		inputbuffer := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := inputbuffer.ReadString('\n')
		fmt.Fprint(connection, text+"\n")
		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		println("Message from server: " + message)
	}
}
