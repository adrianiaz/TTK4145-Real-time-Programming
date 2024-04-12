package orderassigner

import (
	. "Project/config"
	. "Project/datatypes"
)

func buttonEventHandler(
	buttonEvent ButtonEvent,
	hallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	cabOrders [NUM_FLOORS]OrderType,
) ([NUM_FLOORS][NUM_HALL_BUTTONS]OrderType, [NUM_FLOORS]OrderType) {

	switch buttonEvent.Button {
	case BT_HallUp, BT_HallDown:
		if hallOrders[buttonEvent.Floor][buttonEvent.Button] != ConfirmedOrder {
			hallOrders[buttonEvent.Floor][buttonEvent.Button] = UnconfirmedOrder
		}
	case BT_Cab:
		cabOrders[buttonEvent.Floor] = ConfirmedOrder
	}

	return hallOrders, cabOrders
}
