package orderassigner

import (
	. "Project/config"
	. "Project/datatypes"
	"Project/elevio"
)

func toggleLights(hallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType, cabOrders [NUM_FLOORS]OrderType) {
	for floor, orders := range hallOrders {
		for button, order := range orders {
			elevio.SetButtonLamp(ButtonType(button), floor, order == ConfirmedOrder)
		}
	}
	for floor, order := range cabOrders {
		elevio.SetButtonLamp(BT_Cab, floor, order == ConfirmedOrder)
	}
}
