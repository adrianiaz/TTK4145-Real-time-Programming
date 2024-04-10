package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	connectionVar := "127.0.0.1:50955"

	addr, err := net.ResolveUDPAddr("udp", connectionVar)
	if err != nil {
		log.Fatal(err)
	}

	listening, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	//message := 0
	var message int
	//fmt.Println(("reciever"))
	for {
		//fmt.Println("newloop")
		buf := make([]byte, 1024)
		//time.Sleep(time.Second * 3)
		err = listening.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, err := listening.Read(buf)
		if err != nil {
			break
		}
		message, _ = strconv.Atoi(string(buf[:n]))
	}

	//close this socket that we are reading from as it is old backup
	listening.Close()

	//new backup
	//linux
	//exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	//mac
	exec.Command("open", "-a", "Terminal", "go", "run", "main.go").Run()

	//primary mode
	//addr, err = net.ResolveUDPAddr("udp", connectionVar)
	//if err != nil {
	//	log.Fatal(err)
	//}

	time.Sleep(1 * time.Second)

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		message++
		_, err := conn.Write([]byte(strconv.Itoa(message)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message)
		time.Sleep(1 * time.Second)
	}

}
