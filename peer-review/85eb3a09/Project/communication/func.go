package communication

import (
	"Project/assigner"
	"Project/config"
	"Project/singleElevator/elevator"
	"Project/singleElevator/elevio"
)

func UpdateAndCompleteOrder(localWorldView *WorldView, button int, floor int, lampStatus bool) {
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, lampStatus)
	localWorldView.CompleteOrder = completeOrderLight{floor, elevio.ButtonType(button), !lampStatus}
}

func ResetAckList(localWorldView *WorldView) {
	localWorldView.AckList = make([]string, 0)
	localWorldView.AckList = append(localWorldView.AckList, localWorldView.ID)
}

func AssignOrder(message WorldView, ch_assignedRequests chan<- map[string][][2]bool) {
	input := HRAInputFormat(message)
	ch_assignedRequests <- assigner.Assigner(input)
}

func InitializeWorldView(elevatorID string) WorldView {
	message := WorldView{
		Counter:         0,
		ID:              elevatorID,
		AckList:         make([]string, 0),
		ElevatorList:    map[string]elevator.Elevator{elevatorID: elevator.MakeElevator()},
		HallOrderStatus: InitializeHallOrderStatus(),
		CompleteOrder:   completeOrderLight{0, elevio.BT_HallUp, true},
	}
	return message
}

func InitializeHallOrderStatus() [][config.NumButtons - 1]config.RequestState { //Sjekk om den har et problem med å sette verdier
	hallOrderStatus := make([][config.NumButtons - 1]config.RequestState, config.NumFloors)
	for floor := range hallOrderStatus {
		for button := range hallOrderStatus[floor] {
			hallOrderStatus[floor][button] = config.None
		}
	}
	return hallOrderStatus
}

func GetHallRequests(message WorldView) [][2]bool {
	hall := make([][2]bool, config.NumFloors)
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < elevio.BT_Cab; button++ {
			if message.HallOrderStatus[floor][button] == config.Confirmed {
				hall[floor][button] = true
			} else {
				hall[floor][button] = false
			}
		}
	}
	return hall
}

func GetCabRequests(e elevator.Elevator) []bool {
	cabRequests := make([]bool, 0)
	for floor := 0; floor < config.NumFloors; floor++ {
		cabRequests = append(cabRequests, e.Requests[floor][elevio.BT_Cab])
	}
	return cabRequests
}

func HRAInputFormat(
	myMessage WorldView,
) assigner.HRAInput {
	elevStates := make(map[string]assigner.HRAElevState)
	hallRequests := GetHallRequests(myMessage)

	for IDs := range myMessage.AckList {
		if !myMessage.ElevatorList[myMessage.AckList[IDs]].Unavailable {
			elevStates[myMessage.AckList[IDs]] = assigner.HRAElevState{
				Behavior:    elevator.StateToString(myMessage.ElevatorList[myMessage.AckList[IDs]].CurrentState),
				Floor:       myMessage.ElevatorList[myMessage.AckList[IDs]].Floor,
				Direction:   elevio.DirToString(myMessage.ElevatorList[myMessage.AckList[IDs]].Dir),
				CabRequests: GetCabRequests(myMessage.ElevatorList[myMessage.AckList[IDs]]),
			}
		}
	}

	input := assigner.HRAInput{
		HallRequests: hallRequests,
		States:       elevStates,
	}
	return input
}