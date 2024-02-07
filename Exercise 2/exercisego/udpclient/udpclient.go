package udpclient

import (
	"fmt"
	"log"
	"net"
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
	fmt.Println(("reciever"))
	for {
		n, err := connection.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("reviever: " + string(buf[:n]))
	}
}

// Try sending a message to the server IP on port 20000 + n
// The server will act the same way if you send a broadcast (#.#.#.255 or 255.255.255.255) instead of
func Udpclient(ipAddr string, port string) {

	/* 	//prompt for user
	   	istream1 := bufio.NewReader(os.Stdin)
	   	fmt.Print("What ip address do you want to connect to? ")
	   	ipAddr, _ := istream1.ReadString('\n')

	   	istream2 := bufio.NewReader(os.Stdin)
	   	fmt.Print("What port do you want to connect to? ")
	   	port, _ := istream2.ReadString('\n')
	*/
	

	//dial up the net
	addr, err := net.ResolveUDPAddr("udp", connectionVar)
	if err != nil {
		log.Fatal(err)
	}
	addr2, err := net.ResolveUDPAddr("udp", ":20005")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	listening, err := net.ListenUDP("udp", addr2)
	if err != nil {
		log.Fatal(err)
	}

	go reciever(listening)
	go transmitter(conn)
	select {}
}
