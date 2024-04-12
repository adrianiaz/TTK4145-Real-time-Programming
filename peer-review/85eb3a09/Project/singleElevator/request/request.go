package request

import (
	"Project/config"
	"Project/singleElevator/elevator"
	"Project/singleElevator/elevio"
	"time"
)

func RequestsAbove(e elevator.Elevator) bool {
	for f := e.Floor + 1; f < config.NumFloors; f++ {
		for btn := 0; btn < len(e.Requests[f]); btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e elevator.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < len(e.Requests[f]); btn++ {
			if e.Requests[f][btn] {
				return true
			}
		}
	}
	return false
}

func RequestClearAtCurrentFloor(e *elevator.Elevator) {
	time.Sleep(50 * time.Millisecond)
	e.Requests[e.Floor][elevio.BT_Cab] = false;
        switch(e.Dir){
        case elevio.MD_Up:
            if(!RequestsAbove(*e) && !e.Requests[e.Floor][elevio.BT_HallUp]){
                e.Requests[e.Floor][elevio.BT_HallDown] = false;
            }
            e.Requests[e.Floor][elevio.BT_HallUp] = false;
        case elevio.MD_Down:
            if(!RequestsBelow(*e) && !e.Requests[e.Floor][elevio.BT_HallDown]){
                e.Requests[e.Floor][elevio.BT_HallUp] = false;
            }
            e.Requests[e.Floor][elevio.BT_HallDown] = false;
        case elevio.MD_Stop:
        default:
            e.Requests[e.Floor][elevio.BT_HallUp] = false;
            e.Requests[e.Floor][elevio.BT_HallDown] = false;
        }
}

func RequestShouldStop(e *elevator.Elevator) bool {
	switch {
	case e.Dir == elevio.MD_Down:
		return e.Requests[e.Floor][int(elevio.BT_HallDown)] || e.Requests[e.Floor][int(elevio.BT_Cab)] || !RequestsBelow(*e)
	case e.Dir == elevio.MD_Up:
		return e.Requests[e.Floor][int(elevio.BT_HallUp)] || e.Requests[e.Floor][int(elevio.BT_Cab)] || !RequestsAbove(*e)
	default:
		return true
	}
}

func RequestChooseDirection(e *elevator.Elevator) {
	switch e.Dir {
	case elevio.MD_Up:
		if RequestsAbove(*e) {
			e.Dir = elevio.MD_Up
		} else if RequestsBelow(*e) {
			e.Dir = elevio.MD_Down
		} else {
			e.Dir = elevio.MD_Stop
		}
	case elevio.MD_Down:
		fallthrough
	case elevio.MD_Stop:
		if RequestsBelow(*e) {
			e.Dir = elevio.MD_Down
		} else if RequestsAbove(*e) {
			e.Dir = elevio.MD_Up
		} else {
			e.Dir = elevio.MD_Stop
		}
	}
}
