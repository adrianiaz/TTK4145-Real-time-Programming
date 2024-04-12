package main

import (
	"assignerwrapper"
	"cab"
	"elevatorglobals"
	"elevatorinterface"
	"flag"
	"fmt"
	"location"
	"motor"
	"network/localip"
	"network/peers"
	"os"
	"strconv"
	"strings"
	"worldviewbuilder"
)

func main() {

	// Program flags

	myElevatorName, elevatorPort := readProgramFlags()
	elevatorglobals.MyElevatorName = myElevatorName

	// Distribution

	elevatorinterface.Init("localhost:" + strconv.Itoa(elevatorPort))

	peerUpdateChannel := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(54320, elevatorglobals.MyElevatorName, peerTxEnable)
	go peers.Receiver(54320, peerUpdateChannel)

	worldviewChannel := make(chan elevatorglobals.Worldview)
	cabStateChannel := make(chan elevatorglobals.CabState)
	orderHandledChannel := make(chan elevatorglobals.OrderEvent)

	go worldviewbuilder.Run(cabStateChannel,
		orderHandledChannel,
		peerUpdateChannel,
		worldviewChannel)

	// Assignment

	assignedOrdersChannel := make(chan elevatorglobals.AssignedOrders)
	go assignerwrapper.Run(worldviewChannel, assignedOrdersChannel)

	// Execution

	startFloor, startBehaviour, startDirection := moveCabToNearestFloorBelow()

	cabState := elevatorglobals.CabState{
		Name:         elevatorglobals.MyElevatorName,
		Behaviour:    startBehaviour,
		Floor:        startFloor,
		Direction:    startDirection,
		MotorWorking: true,
		Obstructed:   false,
	}

	go cab.Run(cabState,
		assignedOrdersChannel,
		cabStateChannel,
		orderHandledChannel)

	fmt.Println("Elevator system initialized")

	for {

	}

}

func moveCabToNearestFloorBelow() (int, elevatorglobals.ElevatorBehaviour, elevatorglobals.Direction) {
	var startFloor int
	startBehaviour := elevatorglobals.ElevatorBehaviour_Idle
	startDirection := elevatorglobals.Direction_Stop

	if location.GetFloor() != -1 {
		startFloor = location.GetFloor()
	} else {
		fmt.Println("main.go: Moving cab to nearest floor below")
		startBehaviour = elevatorglobals.ElevatorBehaviour_Idle
		motor.SetDirection(elevatorglobals.Direction_Down)
		startDirection = elevatorglobals.Direction_Down

		for {
			if location.GetFloor() != -1 {
				motor.SetDirection(elevatorglobals.Direction_Stop)
				startFloor = location.GetFloor()
				break
			}
		}
	}

	return startFloor, startBehaviour, startDirection
}

func readProgramFlags() (string, int) {
	localIP, _ := localip.LocalIP()
	elevatorName := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	flag.StringVar(&elevatorName, "name", elevatorName, "name of elevator (string)")

	elevatorPort := 15657
	flag.IntVar(&elevatorPort, "port", 15657, "Elevator server port (uint16)")

	flag.Parse()

	if !strings.Contains(elevatorName, elevatorglobals.Codeword) {
		fmt.Println("main.go: id must contain codeword to connect to network \"" + elevatorglobals.Codeword + "\"\n")
		return elevatorName, elevatorPort
	}
	return elevatorName, elevatorPort
}
