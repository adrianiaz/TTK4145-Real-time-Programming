package elevatordriver

import (
	. "Project/config"
	. "Project/datatypes"
	"Project/elevio"
	"Project/phoenix"
	"fmt"
	"time"
)

type Controller struct {
	moving                bool
	doorOpen              bool
	sendStatus            bool
	obstruction           bool
	lastKnownMovement     time.Time
	mostRecentDoorOpening time.Time
	elevator              Elevator
	orders                [NUM_FLOORS][NUM_BUTTONS]bool
}

type DirectionBehaviourPair struct {
	direction MotorDirection
	behavior  ElevatorBehaviour
}

func ElevatorDriver(
	fromOrderAssignerChannel <-chan OrderAssignerToElevatorPayload,
	toOrderAssignerChannel chan<- ElevatorToOrderAssignerPayload,
	lifelineChannel chan<- bool,
) {
	defer phoenix.RespawnAfterPanic()

	var (
		controller             Controller = initializeController()
		prevController         Controller = controller
		doorTimeout            time.Time
		obstructionChannel     chan bool = make(chan bool)
		floorSensorChannel     chan int  = make(chan int)
		obstructionQuitChannel chan bool = make(chan bool)
	)

	go elevio.PollFloorSensor(floorSensorChannel)
	go elevio.PollObstructionSwitch(obstructionChannel, obstructionQuitChannel)

	// init and get to known state
	controller.obstruction = <-obstructionChannel
	for controller.obstruction {
		elevio.SetDoorOpenLamp(controller.obstruction)
		controller.obstruction = <-obstructionChannel
	}
	obstructionQuitChannel <- true
	elevio.SetDoorOpenLamp(controller.obstruction)
	elevio.SetMotorDirection(controller.elevator.Direction)
	controller.elevator.Floor = <-floorSensorChannel
	elevio.SetFloorIndicator(controller.elevator.Floor)
	elevio.SetMotorDirection(MD_Stop)
	controller.elevator.Behaviour = EB_DoorOpen
	controller.elevator.Direction = MD_Up
	toOrderAssignerChannel <- ElevatorToOrderAssignerPayload(controller.elevator)

	fmt.Println("Driver initialized")

	for {
		if !selfCheck(controller) {
			fmt.Println("Elevator did not pass self diagnosis.")
			panic(controller)
		}
		lifelineChannel <- true
		prevController = controller

		select {
		case controller.orders = <-fromOrderAssignerChannel:
		case controller.elevator.Floor = <-floorSensorChannel:
			elevio.SetFloorIndicator(controller.elevator.Floor)
			controller.lastKnownMovement = time.Now()
		case controller.obstruction = <-obstructionChannel:
		default:
		}

		switch controller.elevator.Behaviour {
		case EB_Idle:
			controller.elevator = updateDirectionAndBehaviour(controller)
		case EB_DoorOpen:
			if !controller.doorOpen {
				controller.doorOpen = true
				controller.mostRecentDoorOpening = time.Now()
				elevio.SetDoorOpenLamp(controller.doorOpen)
				controller.orders = updateOrdersOnDoorOpen(controller)
				doorTimeout = time.Now().Add(ELEVATOR_DOOR_OPEN_TIME)
				go elevio.PollObstructionSwitch(obstructionChannel, obstructionQuitChannel)
			} else {
				if time.Now().After(doorTimeout) {
					controller.doorOpen = false
					elevio.SetDoorOpenLamp(controller.doorOpen)
					controller.elevator = updateDirectionAndBehaviour(controller)
					obstructionQuitChannel <- true
				}
				if controller.obstruction {
					doorTimeout = time.Now().Add(ELEVATOR_DOOR_OPEN_TIME)
				}
			}
		case EB_Moving:
			if !controller.moving {
				controller.moving = true
				controller.lastKnownMovement = time.Now()
				elevio.SetMotorDirection(controller.elevator.Direction)
			}
			if shouldStop(controller) && controller.elevator.Floor != prevController.elevator.Floor {
				elevio.SetMotorDirection(MD_Stop)
				controller.moving = false
				controller.elevator.Behaviour = EB_DoorOpen
			}

		}
		if controller != prevController {
			controller.sendStatus = true
		}
		if controller.sendStatus {
			select {
			case toOrderAssignerChannel <- ElevatorToOrderAssignerPayload(controller.elevator):
				controller.sendStatus = false
			default:
			}
		}
	}
}
