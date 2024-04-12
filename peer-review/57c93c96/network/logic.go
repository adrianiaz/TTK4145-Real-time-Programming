package network

import (
	. "Project/config"
	. "Project/datatypes"
)

func cyclicCounterLogicNetwork(
	hallOrderList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	aliveList [NUM_ELEVATORS]bool,
	myID int,
) [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType {

	for floor := range hallOrderList[0] {
		for button := range hallOrderList[0][0] {
			transition := true
			myOrder := hallOrderList[myID][floor][button]
			for node := range hallOrderList {
				nodeOrder := hallOrderList[node][floor][button]

				if !aliveList[node] || node == myID || nodeOrder == Garbage {
					continue
				}

				if nodeOrder != UnconfirmedOrder {
					transition = false
				}

				myOrder = cyclicHelper(myOrder, nodeOrder)
			}

			if transition && myOrder == UnconfirmedOrder {
				myOrder = ConfirmedOrder
			}
			hallOrderList[myID][floor][button] = myOrder
		}
	}
	return hallOrderList[myID]
}

func cyclicHelper(myOrder OrderType, nodeOrder OrderType) OrderType {
	if myOrder == Garbage {
		myOrder = nodeOrder
	} else if myOrder == NoOrder && nodeOrder == UnconfirmedOrder {
		myOrder = UnconfirmedOrder
	} else if myOrder == UnconfirmedOrder && nodeOrder == ConfirmedOrder {
		myOrder = ConfirmedOrder
	} else if myOrder == ConfirmedOrder && nodeOrder == NoOrder {
		myOrder = NoOrder
	}
	return myOrder
}

func cabOrderOverwriteGarbage(
	currentCabs [NUM_ELEVATORS][NUM_FLOORS]OrderType,
	previousCabs [NUM_ELEVATORS][NUM_FLOORS]OrderType,
) [NUM_ELEVATORS][NUM_FLOORS]OrderType {

	for elevator := range currentCabs {
		for floor, value := range currentCabs[elevator] {
			if value == Garbage {
				currentCabs[elevator][floor] = previousCabs[elevator][floor]
			}
		}
	}
	return currentCabs
}

func cabOrderUpdateFromAssigner(
	cabFromAssigner [NUM_FLOORS]OrderType,
	cabFromNetwork [NUM_FLOORS]OrderType,
) [NUM_FLOORS]OrderType {
	for floor, value := range cabFromAssigner {
		if value != Garbage {
			cabFromNetwork[floor] = cabFromAssigner[floor]
		}
	}
	return cabFromNetwork
}

func elevatorEvent(myID int,
	aliveList [NUM_ELEVATORS]bool,
	elevator Elevator,
	hallOrdersList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
) [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType {

	if elevator.Behaviour != EB_DoorOpen {
		return hallOrdersList[myID]
	}

	floor := elevator.Floor

	for button := range hallOrdersList[0][0] {
		transition := true
		for node := range hallOrdersList {
			if hallOrdersList[node][floor][button] != ConfirmedOrder && aliveList[node] {
				transition = false
			}
		}
		if transition {
			if elevator.Direction == MD_Stop {
				hallOrdersList[myID][floor][button] = NoOrder
			} else if button == int(BT_HallUp) && elevator.Direction == MD_Up {
				hallOrdersList[myID][floor][button] = NoOrder
			} else if button == int(BT_HallDown) && elevator.Direction == MD_Down {
				hallOrdersList[myID][floor][button] = NoOrder
			}
		}
	}
	return hallOrdersList[myID]
}

func buttonEvent(
	hallOrderList [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	messageFromAssigner [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
) [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType {

	for floor := range hallOrderList {
		for button, order := range hallOrderList[floor] {
			if order != ConfirmedOrder && messageFromAssigner[floor][button] == UnconfirmedOrder {
				hallOrderList[floor][button] = UnconfirmedOrder
			}
		}
	}
	return hallOrderList
}
