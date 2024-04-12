package orderassigner

import (
	. "Project/config"
	. "Project/datatypes"
)

func payloadToNetworkInit(elevator Elevator) OrderAssignerToNetworkPayload {

	var hallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType
	var caborders [NUM_FLOORS]OrderType
	for floor := range hallOrders {
		for button := range hallOrders[floor] {
			hallOrders[floor][button] = Garbage
		}
		caborders[floor] = Garbage
	}

	payload := OrderAssignerToNetworkPayload{
		HallOrders: hallOrders,
		CabOrders:  caborders,
		Elevator:   elevator,
	}
	return payload
}
