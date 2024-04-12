package datatypes

import (
	. "Project/config"
)

type (
	ThreadID                       int
	OrderType                      int
	ButtonType                     int
	MotorDirection                 int
	ElevatorBehaviour              int
	ElevatorToOrderAssignerPayload Elevator
	OrderAssignerToElevatorPayload [NUM_FLOORS][NUM_BUTTONS]bool

	ButtonEvent struct {
		Floor  int
		Button ButtonType
	}

	Elevator struct {
		Floor     int
		Direction MotorDirection
		Behaviour ElevatorBehaviour
	}

	NetworkToOrderAssignerPayload struct {
		HallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType
		CabOrders  [NUM_ELEVATORS][NUM_FLOORS]OrderType
		Elevators  [NUM_ELEVATORS]Elevator
		Alive      [NUM_ELEVATORS]bool
	}

	OrderAssignerToNetworkPayload struct {
		HallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType
		CabOrders  [NUM_FLOORS]OrderType
		Elevator   Elevator
	}
)

const (
	MD_Up   MotorDirection = 1
	MD_Down MotorDirection = -1
	MD_Stop MotorDirection = 0

	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2

	EB_Idle     ElevatorBehaviour = 0
	EB_DoorOpen ElevatorBehaviour = 1
	EB_Moving   ElevatorBehaviour = 2

	Garbage          OrderType = -1
	NoOrder          OrderType = 0
	UnconfirmedOrder OrderType = 1
	ConfirmedOrder   OrderType = 2
)
