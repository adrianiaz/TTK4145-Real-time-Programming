package communication

import (
	"Project/config"
	"Project/network/peers"
	"Project/saveJSON"
	"Project/singleElevator/elevator"
	"Project/singleElevator/elevio"
	"Project/singleElevator/fsm"
	"fmt"
	"reflect"
	"time"
)

type WorldView struct {
	Counter         int
	ID              string
	AckList         []string
	ElevatorList    map[string]elevator.Elevator
	HallOrderStatus [][config.NumButtons - 1]config.RequestState
	CompleteOrder   completeOrderLight
}

type completeOrderLight struct {
	Floor    int
	Button   elevio.ButtonType
	LightOff bool
}

func CommunicationHandler(
	elevatorID string,
	ch_CabOrders chan<- elevio.ButtonEvent,
	ch_assignedRequests chan<- map[string][][2]bool,
	ch_messageTx chan<- WorldView,
	ch_messageRx <-chan WorldView,
	ch_newLocalElevator <-chan elevator.Elevator,
	ch_peerUpdate <-chan peers.PeerUpdate,
	ch_order <-chan elevio.ButtonEvent,
	ch_completeOrder <-chan fsm.DoneOrder,
	ch_confirmedOrder chan WorldView,
	ch_mergeFSM chan<- elevio.ButtonEvent,
	ch_peerTXEnable chan<- bool,
) {
	//Initilizing local elevator message
	initlocalWorldView := InitializeWorldView(elevatorID)
	localWorldView := &initlocalWorldView

	//Defining timer
	SendMessageTimer := time.NewTimer(time.Duration(config.SendMessageTimerMs) * time.Millisecond)

	//Defining variables
	numPeers := 0
	orderDistributed := make([][config.NumButtons - 1]bool, config.NumFloors)

	for {
	OuterLoop:
		select {
		case order := <-ch_order: 

			if numPeers == 1 { 
				ch_mergeFSM <- order 
				UpdateAndCompleteOrder(localWorldView, int(order.Button), order.Floor, true)
				break OuterLoop
			}
			if order.Button == elevio.BT_Cab {
				ch_CabOrders <- order 
				break OuterLoop
			}
			localWorldView.HallOrderStatus[order.Floor][int(order.Button)] = config.Order 
			localWorldView.Counter++                                                      
			ResetAckList(localWorldView)

		case updatedWorldView := <-ch_messageRx: 

			if localWorldView.Counter >= updatedWorldView.Counter { 
				if localWorldView.Counter == updatedWorldView.Counter && len(localWorldView.AckList) < len(updatedWorldView.AckList) {
					//Updating local world view with world view recieved from messageRx, without overwriting the local elevator in the elevator list
					localElevatorStatus := localWorldView.ElevatorList[elevatorID]
					localWorldView = &updatedWorldView
					localWorldView.ElevatorList[elevatorID] = localElevatorStatus
				} else {
					break OuterLoop
				}
			}

			if updatedWorldView.CompleteOrder.LightOff {
				elevio.SetButtonLamp(updatedWorldView.CompleteOrder.Button, updatedWorldView.CompleteOrder.Floor, false)
				localWorldView.CompleteOrder.LightOff = true
			}

			if len(updatedWorldView.AckList) == numPeers { 
				for floor := 0; floor < config.NumFloors; floor++ {
					for button := 0; button < config.NumButtons-1; button++ { 
						switch {
						case updatedWorldView.HallOrderStatus[floor][button] == config.Order:
							localWorldView.HallOrderStatus[floor][button] = config.Confirmed 
							localWorldView.Counter = updatedWorldView.Counter                
							localWorldView.Counter++                                         
							ResetAckList(localWorldView)
						case updatedWorldView.HallOrderStatus[floor][button] == config.Confirmed && !orderDistributed[floor][button]:
							UpdateAndCompleteOrder(localWorldView, button, floor, true)
							AssignOrder(updatedWorldView, ch_assignedRequests)
							orderDistributed[floor][button] = true 
							localWorldView = &updatedWorldView
							localWorldView.ID = elevatorID
						case updatedWorldView.HallOrderStatus[floor][button] == config.Complete:
							localWorldView.HallOrderStatus[floor][button] = config.None
							orderDistributed[floor][button] = false
							localWorldView.Counter++
						}
					}
				}

			} else {
				for IDs := range updatedWorldView.AckList { 
					if localWorldView.ID == updatedWorldView.AckList[IDs] {
					
						if reflect.DeepEqual(localWorldView.AckList, updatedWorldView.AckList) { 
							//Updating local world view with world view recieved from messageRx, without overwriting the local elevator in the elevator list
							localElevatorStatus := localWorldView.ElevatorList[elevatorID]
							localWorldView = &updatedWorldView
							localWorldView.ElevatorList[elevatorID] = localElevatorStatus

							break OuterLoop
						}
						//Updating local world view with world view recieved from messageRx, without overwriting the local elevator in the elevator list
						localElevatorStatus := localWorldView.ElevatorList[elevatorID]
						localWorldView = &updatedWorldView
						localWorldView.ElevatorList[elevatorID] = localElevatorStatus

						localWorldView.Counter++
						break OuterLoop
					}
				}
				//Updating local world view with world view recieved from messageRx, without overwriting the local elevator in the elevator list
				localElevatorStatus := localWorldView.ElevatorList[elevatorID]
				localWorldView = &updatedWorldView
				localWorldView.ElevatorList[elevatorID] = localElevatorStatus

				
				localWorldView.AckList = append(localWorldView.AckList, elevatorID) 
				localWorldView.Counter++

				if len(updatedWorldView.AckList) == numPeers { 
					for floor := 0; floor < config.NumFloors; floor++ {
						for button := 0; button < config.NumButtons-1; button++ { 
							if localWorldView.HallOrderStatus[floor][button] == config.Confirmed && !orderDistributed[floor][button] { 
								UpdateAndCompleteOrder(localWorldView, button, floor, true)
								AssignOrder(updatedWorldView, ch_assignedRequests)
								orderDistributed[floor][button] = true
							}

						}
					}
				}
			}

		case <-SendMessageTimer.C:
			localWorldView.ID = elevatorID 
			ch_messageTx <- *localWorldView
			SendMessageTimer.Reset(time.Duration(config.SendMessageTimerMs) * time.Millisecond)

		case complete := <-ch_completeOrder:
			for button := 0; button < len(complete.Order)-1; button++ {
				if complete.Order[button] { 
					localWorldView.HallOrderStatus[complete.Floor][button] = config.Complete
					UpdateAndCompleteOrder(localWorldView, button, complete.Floor, false)
					orderDistributed[complete.Floor][button] = false
				}
			}
			localWorldView.Counter++

		case elev := <-ch_newLocalElevator:
			localWorldView.ElevatorList[elevatorID] = elev
			cabReq := GetCabRequests(elev)
			saveJSON.SaveCabButtonToFile(cabReq, "cabOrders.json")

		case peers := <-ch_peerUpdate:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", peers.Peers)
			fmt.Printf("  New:      %q\n", peers.New)
			fmt.Printf("  Lost:     %q\n", peers.Lost)
			
			numPeers = len(peers.Peers)

			if len(peers.Lost) > 0 {
				if localWorldView.ElevatorList[peers.Lost[0]].Unavailable {
					AssignOrder(*localWorldView, ch_assignedRequests)
					ch_peerTXEnable <- true
				} else {
					for i, ack := range localWorldView.AckList {
						for _, lostPeer := range peers.Lost {
							delete(localWorldView.ElevatorList, lostPeer)
							if ack == lostPeer {
								localWorldView.AckList = append(localWorldView.AckList[:i], localWorldView.AckList[i+1:]...)
							}
						}
					}
					AssignOrder(*localWorldView, ch_assignedRequests)
				}
			}
		}
	}
}
