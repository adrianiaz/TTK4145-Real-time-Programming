package elevator

import (
	"Project/config"
	"Project/singleElevator/elevio"
)

type State int

const (
	Idle     State = 0
	DoorOpen State = 1
	Moving   State = 2
)

type Elevator struct {
	Floor        int
	Dir          elevio.MotorDirection
	Requests     [][]bool
	CurrentState State
	TimerCount   int
	Unavailable  bool
}

// Initialize the local elevator without orders and in floor 0 in IDLE state.
func MakeElevator() Elevator {
	requests := make([][]bool, 0)
	for floor := 0; floor < config.NumFloors; floor++ {
		requests = append(requests, make([]bool, config.NumButtons))
		for button := 0; button < config.NumButtons; button++ {
			requests[floor][button] = false
		}
	}
	return Elevator{
		Floor:        0,
		Dir:          elevio.MD_Stop,
		Requests:     requests,
		CurrentState: Idle,
		TimerCount:   0,
		Unavailable:  false}
}

// Check for cab orders and updates the lights accordingly.
func SetLamps(e Elevator) {
	elevio.SetFloorIndicator(e.Floor)
	for f := 0; f < config.NumFloors; f++ {
		elevio.SetButtonLamp(elevio.ButtonType(elevio.BT_Cab), f, e.Requests[f][elevio.BT_Cab])
	}
}

func StateToString(s State) string {
	switch s {
	case Idle:
		return "idle"
	case DoorOpen:
		return "doorOpen"
	case Moving:
		return "moving"
	default:
		return "undefined"
	}
}
