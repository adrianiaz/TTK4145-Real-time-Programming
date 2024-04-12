package orderassigner

import (
	. "Project/datatypes"
	"Project/elevio"
	"fmt"
)

func OrderAssigner(
	fromNetworkChannel <-chan NetworkToOrderAssignerPayload,
	toNetworkChannel chan<- OrderAssignerToNetworkPayload,
	fromElevatorChannel <-chan ElevatorToOrderAssignerPayload,
	toElevatorChannel chan<- OrderAssignerToElevatorPayload,
	lifelineChannel chan<- bool,
	myID int,
) {

	var (
		sendToNetwork        bool             = false
		sendToElevator       bool             = false
		receivedFromButton   bool             = false
		receivedFromElevator bool             = false
		receivedFromNetwork  bool             = false
		buttonPollerChannel  chan ButtonEvent = make(chan ButtonEvent)

		payloadToElevator   OrderAssignerToElevatorPayload
		payloadFromElevator ElevatorToOrderAssignerPayload
		payloadToNetwork    OrderAssignerToNetworkPayload
		payloadFromNetwork  NetworkToOrderAssignerPayload
		buttonEvent         ButtonEvent
	)
	// init and get to known state
	payloadFromElevator = <-fromElevatorChannel
	payloadToNetwork = payloadToNetworkInit(Elevator(payloadFromElevator))
	sendToNetwork = true
	toggleLights(payloadToNetwork.HallOrders, payloadToNetwork.CabOrders)

	go elevio.PollButtons(buttonPollerChannel)

	fmt.Println("Assigner initialized")

	for {
		select {
		case lifelineChannel <- true:
		case buttonEvent = <-buttonPollerChannel:
			receivedFromButton = true
		case payloadFromElevator = <-fromElevatorChannel:
			receivedFromElevator = true
		case payloadFromNetwork = <-fromNetworkChannel:
			receivedFromNetwork = true
		}

		if receivedFromButton {
			payloadToNetwork.HallOrders, payloadToNetwork.CabOrders = buttonEventHandler(
				buttonEvent,
				payloadToNetwork.HallOrders,
				payloadToNetwork.CabOrders,
			)
			sendToNetwork = true

			receivedFromButton = false
		}

		if receivedFromElevator {
			if payloadFromElevator.Behaviour == EB_DoorOpen {
				payloadToNetwork.CabOrders[payloadFromElevator.Floor] = NoOrder
			}

			payloadToNetwork.Elevator = Elevator(payloadFromElevator)
			sendToNetwork = true

			receivedFromElevator = false
		}

		if receivedFromNetwork {
			toggleLights(payloadFromNetwork.HallOrders, payloadFromNetwork.CabOrders[myID])
			payloadToNetwork.HallOrders = payloadFromNetwork.HallOrders

			payloadToNetwork.CabOrders = cabOrderFilterFromNetwork(
				payloadToNetwork.CabOrders,
				payloadFromNetwork.CabOrders[myID],
			)

			payloadToElevator = CostFunction(
				payloadFromNetwork.HallOrders,
				orderTypeToBool(payloadFromNetwork.CabOrders),
				payloadFromNetwork.Elevators,
				myID,
				payloadFromNetwork.Alive,
			)
			sendToElevator = true

			receivedFromNetwork = false
		}

		if sendToNetwork {
			toNetworkChannel <- payloadToNetwork
			sendToNetwork = false
		}

		if sendToElevator {
			toElevatorChannel <- payloadToElevator
			sendToElevator = false
		}
	}
}
