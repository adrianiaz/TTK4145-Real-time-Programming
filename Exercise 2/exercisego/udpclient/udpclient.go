package udpclient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func transmitter(connection *net.UDPConn) {
	for {
		_, err := connection.Write([]byte("Heeellloooooo!"))
		if err != nil {
			log.Fatal(err)
		}
		timer := time.NewTimer(1 * time.Second) //send message and wait 1 second
		<-timer.C
	}
}
func reciever(connection *net.UDPConn) {
	buf := make([]byte, 1024)
	for {

		n, err := connection.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buf[:n]))

	}
}

// Try sending a message to the server IP on port 20000 + n
// The server will act the same way if you send a broadcast (#.#.#.255 or 255.255.255.255) instead of
func Udpclient() {

	//prompt for user
	istream1 := bufio.NewReader(os.Stdin)
	fmt.Print("What ip address do you want to connect to? ")
	ipAddr, _ := istream1.ReadString('\n')

	istream2 := bufio.NewReader(os.Stdin)
	fmt.Print("What port do you want to connect to? ")
	port, _ := istream2.ReadString('\n')

	//dial up the net
	addr, err := net.ResolveUDPAddr("tcp4", ipAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	go transmitter(conn)
	go reciever(conn)
}
