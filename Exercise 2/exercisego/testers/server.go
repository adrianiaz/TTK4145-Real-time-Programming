package testers

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func Server() {
	//create socket
	socket, err := net.Listen("tcp4", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listening on socket 8000")
	//accept connections to the listening socket
	conn, err := socket.Accept()
	if err != nil {
		log.Fatal(err)
	}

	//continiously look for message from clients, and confirm that message is recieved
	for {
		//close connection if no message is recieved in 5 seconds
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message from client: ", message)

		var returnMessagetoClient string = "This is the message we recieved from you: "
		conn.Write([]byte(returnMessagetoClient + message + "\n"))
	}

}
