package tcpclient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

//Fixed size messages of size 1024, if you connect to port 34933
//Delimited messages that use \0 as the marker, if you connect to port 33546
//Ip for computer at lab is 10.100.23.186

//Tell the server to connect back to you, by sending a message of the form Connect to: #.#.#.#:#\0
//(IP of your machine and port you are listening to). You can find your own address by running ifconfig in the terminal

func Tcpclient(ipAddr string, port string) {
	/* istream1 := bufio.NewReader(os.Stdin)
	fmt.Print("What ip address do you want to connect to? ")
	ipAddr, _ := istream1.ReadString(':')

	istream2 := bufio.NewReader(os.Stdin)
	fmt.Print("What port do you want to connect to? ")
	port, _ := istream2.ReadString('\n') */

	connectionVar := ipAddr + ":" + port
	fmt.Println(connectionVar)

	connection, err := net.Dial("tcp", connectionVar)
	if err != nil {
		log.Fatal(err)
	}

	for {

		//writing to server
		inputbuffer := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send to server (\"quit\" to close connection): ")
		text, _ := inputbuffer.ReadString('\n')
		if text == "quit\n" {
			connection.Close()
			return
		}
		if port == "33546" { // this port uses null delimiter
			fmt.Fprint(connection, text+"\x00")
		} else {
			fmt.Fprint(connection, text)
		}

		//Recieving message from server
		reply, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		println("Message from server: \"" + reply + "\"")
	}

}
