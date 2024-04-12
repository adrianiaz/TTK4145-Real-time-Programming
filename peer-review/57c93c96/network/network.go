package network

import (
	. "Project/config"
	. "Project/datatypes"
	"fmt"
)

type PayloadOnNetwork struct {
	Id         int
	HallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType
	CabOrders  [NUM_ELEVATORS][NUM_FLOORS]OrderType
	Elevator   Elevator
}

func NetworkHandler(
	fromOrderAssignerChannel <-chan OrderAssignerToNetworkPayload,
	toOrderAssignerChannel chan<- NetworkToOrderAssignerPayload,
	lifelineChannel chan<- bool,
	myID int,
) {

	var (
		sendToAssigner         bool = false
		receivedFromAssigner   bool = false
		receivedFromNetwork    bool = false
		connectedToNetwork     bool = false
		prevConnectetToNetwork bool = false

		aliveList     [NUM_ELEVATORS]bool
		elevatorList  [NUM_ELEVATORS]Elevator
		cabOrderList  [NUM_ELEVATORS][NUM_FLOORS]OrderType                   = cabOrderListInit()
		hallOrderList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType = hallOrderListInit()

		messageFromNetwork  PayloadOnNetwork
		messageFromAssigner OrderAssignerToNetworkPayload

		fromReceiverChannel     chan PayloadOnNetwork    = make(chan PayloadOnNetwork)
		toTransmitterChannel    chan PayloadOnNetwork    = make(chan PayloadOnNetwork)
		connectionStatusChannel chan [NUM_ELEVATORS]bool = make(chan [NUM_ELEVATORS]bool)
	)
	// init and get to known state
	messageFromAssigner = <-fromOrderAssignerChannel
	receivedFromAssigner = true

	go receiver(fromReceiverChannel, connectionStatusChannel, myID)
	go transmitter(toTransmitterChannel)

	fmt.Println("Network initialized")

	for {
		lifelineChannel <- true
		prevConnectetToNetwork = connectedToNetwork

		select {
		case messageFromAssigner = <-fromOrderAssignerChannel:
			receivedFromAssigner = true
		case messageFromNetwork = <-fromReceiverChannel:
			receivedFromNetwork = true
		case aliveList = <-connectionStatusChannel:
			connectedToNetwork = aliveList[myID]
			sendToAssigner = true
		default:
		}

		if !prevConnectetToNetwork && connectedToNetwork {
			hallOrderList = hallOrderListReset(myID, hallOrderList)
			cabOrderList = cabOrderListReset(myID, cabOrderList)
		}

		if receivedFromAssigner {
			elevatorList[myID] = messageFromAssigner.Elevator
			cabOrderList[myID], hallOrderList[myID] = handleInputFromAssigner(myID, messageFromAssigner, hallOrderList, cabOrderList[myID], aliveList, connectedToNetwork)
			if !connectedToNetwork {
				sendToAssigner = true
			}
			receivedFromAssigner = false
		}

		if receivedFromNetwork {
			elevatorList[messageFromNetwork.Id] = messageFromNetwork.Elevator
			hallOrderList, cabOrderList, sendToAssigner = handleInputFromNetwork(myID, messageFromNetwork, hallOrderList, cabOrderList, aliveList)
			receivedFromNetwork = false
		}

		if sendToAssigner {
			select {
			case toOrderAssignerChannel <- NetworkToOrderAssignerPayload{
				HallOrders: hallOrderList[myID],
				CabOrders:  cabOrderList,
				Elevators:  elevatorList,
				Alive:      aliveList,
			}:
				sendToAssigner = false
			default:
			}
		}

		if connectedToNetwork {
			select {
			case toTransmitterChannel <- PayloadOnNetwork{
				Id:         myID,
				HallOrders: hallOrderList[myID],
				CabOrders:  cabOrderList,
				Elevator:   elevatorList[myID],
			}:
			default:
			}
		} else {
			select {
			case toTransmitterChannel <- PayloadOnNetwork{
				Id:         myID,
				HallOrders: hallOrderListReset(myID, hallOrderList)[myID],
				CabOrders:  cabOrderListReset(myID, cabOrderList),
				Elevator:   elevatorList[myID],
			}:
			default:
			}
		}
	}
}
