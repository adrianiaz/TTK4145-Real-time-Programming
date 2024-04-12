package network

import (
	. "Project/config"
	. "Project/datatypes"
)

func hallOrderListReset(  
	myID int,
	currentHallOrders [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
) [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType {

	var hallOrderList = hallOrderListInit()
	for floor := range currentHallOrders[myID] {
		for button, value := range currentHallOrders[myID][floor] {
			if value == ConfirmedOrder {
				hallOrderList[myID][floor][button] = ConfirmedOrder
			}
		}
	}
	return hallOrderList
}

func hallOrderListInit() [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType {
	var hallOrderList [NUM_ELEVATORS][NUM_FLOORS][NUM_HALL_BUTTONS]OrderType
	for node := range hallOrderList {
		for floor := range hallOrderList[node] {
			for button := range hallOrderList[node][floor] {
				hallOrderList[node][floor][button] = Garbage
			}
		}
	}
	return hallOrderList
}

func cabOrderListReset(
	myID int,
	currentCabOrderList [NUM_ELEVATORS][NUM_FLOORS]OrderType,
) [NUM_ELEVATORS][NUM_FLOORS]OrderType {
	var cabOrderList = cabOrderListInit()
	for floor := range cabOrderList[myID] {
		cabOrderList[myID][floor] = currentCabOrderList[myID][floor]
	}
	return cabOrderList
} 

func cabOrderListInit() [NUM_ELEVATORS][NUM_FLOORS]OrderType {
	var cabOrderList [NUM_ELEVATORS][NUM_FLOORS]OrderType
	for node := range cabOrderList {
		for floor := range cabOrderList[node] {
			cabOrderList[node][floor] = Garbage
		}
	}
	return cabOrderList
}