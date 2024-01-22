package testers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

//Fixed size messages of size 1024, if you connect to port 34933
//Delimited messages that use \0 as the marker, if you connect to port 33546
//Ip for computer at lab is 10.100.23.129

//Tell the server to connect back to you, by sending a message of the form Connect to: #.#.#.#:#\0
//(IP of your machine and port you are listening to). You can find your own address by running ifconfig in the terminal

func Client() {
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
