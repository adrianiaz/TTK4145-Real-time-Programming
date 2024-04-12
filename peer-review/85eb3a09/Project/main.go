package main

import (
	"Project/communication"
	"Project/config"
	"Project/helpers"
	"Project/network/bcast"
	"Project/network/peers"
	"Project/saveJSON"
	"Project/singleElevator/elevator"
	"Project/singleElevator/elevio"
	"Project/singleElevator/fsm"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
)

func main() {

	//Set up flags
	var id string
	var port string
	var portBackup string
	var backup string
	
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.StringVar(&port, "port", "", "port number to use for communication")
	flag.StringVar(&portBackup, "portBackup", "", "portBackup of this peer")
	flag.StringVar(&backup, "backup", "", "backup")
	flag.Parse()

	//Initialize the elevator on the network
	elevio.Init("localhost:"+port, config.NumFloors)

	//Channels for CommunicationHandler and FSM
	ch_orders := make(chan elevio.ButtonEvent, 100)
	ch_arrivedAtfloors := make(chan int)
	ch_obstruction := make(chan bool)
	ch_localElevatorState := make(chan elevator.Elevator, 100)
	ch_timerDoor := make(chan bool)
	ch_localElevatorRequest := make(chan map[string][][2]bool, 100)
	ch_merged := make(chan elevio.ButtonEvent, 100)
	ch_orderCompleted := make(chan fsm.DoneOrder, 100)
	ch_CabOrders := make(chan elevio.ButtonEvent, 100)
	ch_confirmedOrder := make(chan communication.WorldView, 100)

	addr, _ := net.ResolveUDPAddr("udp", "localhost:"+portBackup)
	helpers.SetUpBackup(backup, ch_CabOrders, addr)

	//////FOR MAC////////
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cmd := exec.Command("osascript", "-e", `tell app "Terminal" to do script "cd `+dir+`&& go run main.go --id=`+id+` --port=`+port+` --backup=1 --portBackup=`+portBackup+`"`)
	// err = cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	/////////////////////

	//////FOR LINUX////////
	fmt.Println("Åpner ny terminal")
	exec.Command("gnome-terminal", "--", "go", "run", "main.go", "--id="+id, "--port="+port, "--backup=1", "--portBackup="+portBackup).Run()
	///////////////////////

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	go helpers.PrimaryAlive(conn)

	//Drivere
	go elevio.PollButtons(ch_orders)
	go elevio.PollFloorSensor(ch_arrivedAtfloors)
	go elevio.PollObstructionSwitch(ch_obstruction)

	//Broadcast channels
	ch_messageTx := make(chan communication.WorldView)
	ch_messageRx := make(chan communication.WorldView)

	go bcast.Transmitter(config.P_BCAST, ch_messageTx)
	go bcast.Receiver(config.P_BCAST, ch_messageRx)

	//Peer to peer channels
	ch_peerUpdate := make(chan peers.PeerUpdate)
	ch_peerTXEnable := make(chan bool)

	go peers.Transmitter(config.P_PEERS, id, ch_peerTXEnable)
	go peers.Receiver(config.P_PEERS, ch_peerUpdate)

	//Communication channels
	saveJSON.RestoreCabOrders(ch_CabOrders, "cabOrders.json")

	go communication.CommunicationHandler(
		id,
		ch_CabOrders,
		ch_localElevatorRequest,
		ch_messageTx,
		ch_messageRx,
		ch_localElevatorState,
		ch_peerUpdate,
		ch_orders,
		ch_orderCompleted,
		ch_confirmedOrder,
		ch_merged,
		ch_peerTXEnable,
	)

	go helpers.MergeChannels(ch_CabOrders, ch_localElevatorRequest, ch_merged, id)

	go fsm.Fsm(
		ch_merged,
		ch_localElevatorState,
		ch_arrivedAtfloors,
		ch_obstruction,
		ch_timerDoor,
		ch_orderCompleted,
		ch_peerTXEnable,
	)

	select {}
}
