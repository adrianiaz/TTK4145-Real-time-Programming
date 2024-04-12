package fsm

import (
	"Project/config"
	"Project/singleElevator/elevator"
	"Project/singleElevator/elevio"
	"Project/singleElevator/request"
	"time"
)

type DoneOrder struct {
	Order []bool
	Floor int
}

type DoneCabOrder struct {
	Order bool
	Floor int
}

// Final state machine that runs the local elevator.
func Fsm(
	ch_orders <-chan elevio.ButtonEvent,
	ch_elevatorState chan<- elevator.Elevator,
	ch_arrivedAtFloors chan int,
	ch_obstruction chan bool,
	ch_timerDoor chan bool,
	ch_completeOrder chan<- DoneOrder,
	ch_peerTXEnable chan bool,
) {

	//Initilizing elevator
	elev := elevator.MakeElevator()
	e := &elev

	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)

	for {
		floor := <-ch_arrivedAtFloors
		if floor != 0 {
			elevio.SetMotorDirection(elevio.MD_Down)
		} else {
			elevio.SetMotorDirection(elevio.MD_Stop)
			break
		}
	}

	// Initializing timers
	doorTimer := time.NewTimer(time.Duration(config.DoorOpenDuration) * time.Second)
	timerUpdateState := time.NewTimer(time.Duration(config.StateUpdatePeriodMs) * time.Microsecond)
	motorStopTimer := time.NewTimer(time.Duration(config.MotorStopDuration) * time.Second)
	motorStopTimer.Stop()

	doneCabOrder := DoneCabOrder{Order: false}

	for {
		elevator.SetLamps(*e)
		select {
		case order := <-ch_orders:
			e.Requests[order.Floor][int(order.Button)] = true
			switch {
			case e.CurrentState == elevator.DoorOpen:
				if e.Floor == order.Floor {
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
					if order.Button == elevio.BT_Cab {
						doneCabOrder.Order = true
						doneCabOrder.Floor = e.Floor
					} else {
						ch_completeOrder <- DoneOrder{Order: e.Requests[e.Floor], Floor: e.Floor}
						request.RequestClearAtCurrentFloor(e)
						ch_elevatorState <- *e
					}
					motorStopTimer.Stop()
				} else {
					e.Requests[order.Floor][int(order.Button)] = true
					motorStopTimer.Reset(time.Duration(config.MotorStopDuration+config.DoorOpenDuration) * time.Second)
				}

			case e.CurrentState == elevator.Moving:
				e.Requests[order.Floor][int(order.Button)] = true

			case e.CurrentState == elevator.Idle:
				motorStopTimer.Stop()
				if e.Floor == order.Floor {
					elevator.SetLamps(*e)
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
					e.CurrentState = elevator.DoorOpen
					if order.Button == elevio.BT_Cab {
						doneCabOrder.Order = true
						doneCabOrder.Floor = e.Floor
					}
					ch_completeOrder <- DoneOrder{Order: e.Requests[e.Floor], Floor: e.Floor}
					ch_elevatorState <- *e
					break
				} else {
					e.Requests[order.Floor][int(order.Button)] = true
					request.RequestChooseDirection(e)
					elevio.SetMotorDirection(e.Dir)
					e.CurrentState = elevator.Moving
					motorStopTimer.Reset(time.Duration(config.MotorStopDuration) * time.Second)
					ch_elevatorState <- *e
					break
				}
			}

		case floor := <-ch_arrivedAtFloors:
			motorStopTimer.Reset(time.Duration(config.MotorStopDuration) * time.Second)
			ch_peerTXEnable <- true
			e.Floor = floor
			switch {
			case e.CurrentState == elevator.Moving:
				if request.RequestShouldStop(e) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					motorStopTimer.Stop()
					elevator.SetLamps(*e)
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
					e.CurrentState = elevator.DoorOpen
					currentOrder := DoneOrder{Order: e.Requests[e.Floor], Floor: e.Floor}

					if currentOrder.Order[0] == currentOrder.Order[1] {
						request.RequestClearAtCurrentFloor(e)
						ch_completeOrder <- DoneOrder{Order: []bool{!currentOrder.Order[0], !currentOrder.Order[1], currentOrder.Order[2]}, Floor: e.Floor} //
					} else {
						ch_completeOrder <- DoneOrder{Order: e.Requests[e.Floor], Floor: e.Floor}
						request.RequestClearAtCurrentFloor(e)
					}
					ch_elevatorState <- *e
				}
			}

		case <-doorTimer.C:
			switch {
			case e.CurrentState == elevator.DoorOpen:
				request.RequestChooseDirection(e)
				elevio.SetMotorDirection(e.Dir)
				elevio.SetDoorOpenLamp(false)

				if e.Dir == elevio.MD_Stop {
					e.CurrentState = elevator.Idle
					ch_elevatorState <- *e
					motorStopTimer.Stop()
				} else {
					e.CurrentState = elevator.Moving
					ch_elevatorState <- *e
					motorStopTimer.Reset(time.Duration(config.MotorStopDuration) * time.Second)
				}
			}
			if doneCabOrder.Order {
				elevio.SetDoorOpenLamp(false)
				e.Requests[doneCabOrder.Floor][elevio.BT_Cab] = false
				doneCabOrder.Order = false
			}

		case obstruction := <-ch_obstruction:
			if obstruction {
				e.Unavailable = true
				ch_elevatorState <- *e
				ch_peerTXEnable <- false
			}
			for obstruction {
				select {
				case obstruction = <-ch_obstruction:
					e.Unavailable = false
					ch_elevatorState <- *e
					doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
				default:
					if e.CurrentState == elevator.DoorOpen {
						doorTimer.Reset(time.Duration(config.DoorOpenDuration) * time.Second)
					}
				}
			}

		case <-timerUpdateState.C:
			ch_elevatorState <- *e
			timerUpdateState.Reset(time.Duration(config.StateUpdatePeriodMs) * time.Millisecond)

		case <-motorStopTimer.C:
			e.Unavailable = true
			ch_elevatorState <- *e
			ch_peerTXEnable <- false
		}
	}
}
