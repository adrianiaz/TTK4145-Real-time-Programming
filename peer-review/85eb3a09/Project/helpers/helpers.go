package helpers

import (
	"Project/saveJSON"
	"Project/singleElevator/elevio"
	"fmt"
	"net"
	"time"
)

func MergeChannels(ch_CabOrders <-chan elevio.ButtonEvent, ch_distributedRequest <-chan map[string][][2]bool, ch_merged chan<- elevio.ButtonEvent, ID string) {
	for {
		select {
		case order := <-ch_CabOrders:
			ch_merged <- order
		case dist := <-ch_distributedRequest:
			for floor, buttons := range dist[ID] {
				for button, status := range buttons {
					if status {
						ch_merged <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(button)}
					}
				}
			}
		}
	}
}

func SetUpBackup(backup string, ch_CabOrders chan elevio.ButtonEvent, addr *net.UDPAddr) {
	if backup == "1" {
		fmt.Println("I am backup")
		conn, _ := net.ListenUDP("udp", addr)
		buffer := make([]byte, 1024)
		for {
			conn.SetDeadline(time.Now().Add(2 * time.Second))
			_, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Primary is dead")
				break
			}
		}
		defer conn.Close()
		saveJSON.RestoreCabOrders(ch_CabOrders, "cabOrders.json")
	}

}

func PrimaryAlive(conn *net.UDPConn) {
	for {
		time.Sleep(500 * time.Millisecond)
		conn.Write([]byte("Primary is up"))
	}
}
