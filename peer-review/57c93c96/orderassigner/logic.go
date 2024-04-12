package orderassigner

import (
	. "Project/config"
	. "Project/datatypes"
)

func cabOrderFilterFromNetwork(
	cabFromAssigner [NUM_FLOORS]OrderType,
	cabFromNetwork [NUM_FLOORS]OrderType,
) [NUM_FLOORS]OrderType {
	for floor, value := range cabFromAssigner {
		if value == Garbage {
			cabFromAssigner[floor] = cabFromNetwork[floor]
		}
	}
	return cabFromAssigner
}

func orderTypeToBool(
	cabOrders [NUM_ELEVATORS][NUM_FLOORS]OrderType,
) [NUM_ELEVATORS][NUM_FLOORS]bool {
	var boolCabOrders [NUM_ELEVATORS][NUM_FLOORS]bool
	for elevator := range cabOrders {
		for floor, value := range cabOrders[elevator] {
			boolCabOrders[elevator][floor] = (value == ConfirmedOrder)
		}
	}
	return boolCabOrders
}
