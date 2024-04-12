package network

import (
	. "Project/config"
	. "Project/datatypes"
)

func handleInputFromNetwork(
	myID int,
	messageFromNetwork PayloadOnNetwork,
	hallOrderList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	cabOrderList [NUM_ELEVATORS][NUM_FLOORS]OrderType,
	aliveList [NUM_ELEVATORS]bool,
) ([NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType, [NUM_ELEVATORS][NUM_FLOORS]OrderType, bool) {
	var sendToAssigner bool = false

	hallOrderList[messageFromNetwork.Id] = messageFromNetwork.HallOrders
	if messageFromNetwork.Id != myID {
		hallOrderList[myID] = cyclicCounterLogicNetwork(hallOrderList, aliveList, myID)
		cabOrderList = cabOrderOverwriteGarbage(messageFromNetwork.CabOrders, cabOrderList)
		sendToAssigner = true
	}

	return hallOrderList, cabOrderList, sendToAssigner
}

func handleInputFromAssigner(
	myID int,
	messageFromAssigner OrderAssignerToNetworkPayload,
	hallOrderList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	cabOrder [NUM_FLOORS]OrderType,
	aliveList [NUM_ELEVATORS]bool,
	connectedToNetwork bool,
) ([NUM_FLOORS]OrderType, [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType) {

	cabOrder = cabOrderUpdateFromAssigner(messageFromAssigner.CabOrders, cabOrder)
	if connectedToNetwork {
		hallOrderList[myID] = buttonEvent(hallOrderList[myID], messageFromAssigner.HallOrders)
	}
	hallOrderList[myID] = elevatorEvent(myID, aliveList, messageFromAssigner.Elevator, hallOrderList)

	return cabOrder, hallOrderList[myID]
}
