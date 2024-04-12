package elevatordriver

import (
	. "Project/config"
	. "Project/datatypes"
)

func updateOrdersOnDoorOpen(
	controller Controller,
) [NUM_FLOORS][NUM_BUTTONS]bool {
	floor := controller.elevator.Floor
	switch controller.elevator.Direction {
	case MD_Down:
		controller.orders[floor][BT_HallDown] = false
	case MD_Up:
		controller.orders[floor][BT_HallUp] = false
	case MD_Stop:
		controller.orders[floor][BT_HallDown] = false
		controller.orders[floor][BT_HallUp] = false
	}
	controller.orders[controller.elevator.Floor][BT_Cab] = false
	return controller.orders
}
