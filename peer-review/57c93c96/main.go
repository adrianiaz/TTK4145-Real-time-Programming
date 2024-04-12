package main

import (
	. "Project/config"
	. "Project/datatypes"
	"Project/elevatordriver"
	"Project/elevio"
	"Project/network"
	"Project/orderassigner"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	elevatorDriverThread ThreadID = 0
	networkThread        ThreadID = 1
	orderAssignerThread  ThreadID = 2
)

func main() {
	myID, respawnTimeout := parseArgs()
	if !(0 <= myID && myID <= (NUM_ELEVATORS-1)) {
		fmt.Printf(`ID %d not valid. ID must be in range: [0-%d]`, myID, (NUM_ELEVATORS - 1))
		return
	}
	defer deathBedProcedure(myID, respawnTimeout)
	if respawnTimeout {
		time.Sleep(TIMEOUT_AFTER_RESPAWN)
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	elevio.Init(ELEVATOR_PORT)

	var (
		assignerToNetworkChannel     = make(chan OrderAssignerToNetworkPayload)
		networkToAssignerChannel     = make(chan NetworkToOrderAssignerPayload)
		assignerToElevatorChannel    = make(chan OrderAssignerToElevatorPayload)
		elevatorToAssignerChannel    = make(chan ElevatorToOrderAssignerPayload)
		networkLifelineChannel       = make(chan bool)
		elevatorLifelineChannel      = make(chan bool)
		orderAssignerLifelineChannel = make(chan bool)
	)

	go elevatordriver.ElevatorDriver(
		assignerToElevatorChannel,
		elevatorToAssignerChannel,
		elevatorLifelineChannel,
	)
	go orderassigner.OrderAssigner(
		networkToAssignerChannel,
		assignerToNetworkChannel,
		elevatorToAssignerChannel,
		assignerToElevatorChannel,
		orderAssignerLifelineChannel,
		myID,
	)
	go network.NetworkHandler(
		assignerToNetworkChannel,
		networkToAssignerChannel,
		networkLifelineChannel,
		myID,
	)

	var threadLastCheckIn = [3]time.Time{
		time.Now(),
		time.Now(),
		time.Now(),
	}

	for {
		select {
		case <-elevatorLifelineChannel:
			threadLastCheckIn[elevatorDriverThread] = time.Now()
		case <-networkLifelineChannel:
			threadLastCheckIn[networkThread] = time.Now()
		case <-orderAssignerLifelineChannel:
			threadLastCheckIn[orderAssignerThread] = time.Now()
		case <-signalChannel:
			fmt.Println("Program terminated by OS or user. Respawning.")
			return
		default:
		}

		for thread, lastCheckIn := range threadLastCheckIn {
			if lastCheckIn.Add(DECLARE_THREAD_DEAD_AFTER).Before(time.Now()) {
				fmt.Printf("Thread %d is dead. Terminating program and respawning.\n", thread)
				return
			}
		}
	}
}

func deathBedProcedure(id int, timeout bool) {

	if timeout {
		exec.Command("gnome-terminal", "--", "go", "run", "main.go", fmt.Sprintf("-id=%d", int(id)), "-timeout=true").Run()
	} else {
		exec.Command("gnome-terminal", "--", "go", "run", "main.go", fmt.Sprintf("-id=%d", int(id))).Run()
	}
	os.Exit(1)
}

func parseArgs() (nodeID int, respawnTimeout bool) {
	flag.IntVar(&nodeID, "id", 0, "Node ID")
	flag.BoolVar(&respawnTimeout, "timeout", false, "Timeout after respawn")
	flag.Parse()
	return nodeID, respawnTimeout
}
