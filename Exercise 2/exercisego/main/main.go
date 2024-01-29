package main

import (
	"exercisego/tcpclient"
	//"exercisego/udpclient"
)

//FOR UDP:
// Try sending a message to the server IP on port 20000 + n
// The server will act the same way if you send a broadcast (#.#.#.255 or 255.255.255.255) instead of

//For TCP:
//Fixed size messages of size 1024, if you connect to port 34933
//Delimited messages that use \0 as the marker, if you connect to port 33546
//Ip for computer at lab is 10.100.23.129

//Tell the server to connect back to you, by sending a message of the form Connect to: #.#.#.#:#\0
//(IP of your machine and port you are listening to). You can find your own address by running ifconfig in the terminal

func main() {
	tcpclient.Tcpclient("10.100.23.186", "34933")
	//udpclient.Udpclient("10.100.23.186", "20005")
}
