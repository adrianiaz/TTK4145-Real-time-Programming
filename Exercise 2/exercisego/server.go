package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func server() {
	//create socket
	socket, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listening on socket 80000")
	//accept connections to the listening socket
	conn, err := socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	//continiously look for message from clients, and confirm that message is recieved
	for {
		//close connection if no message is recieved in 5 seconds
		conn.SetDeadline(time.Now().Add(5 * time.Second))

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message from client: ", message)

		var returnMessagetoClient string = "Message recieved"
		conn.Write([]byte(returnMessagetoClient + "\n"))
	}

}
