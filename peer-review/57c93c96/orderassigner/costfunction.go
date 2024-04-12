package orderassigner

import (
	. "Project/config"
	. "Project/datatypes"
	"Project/phoenix"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"runtime"
)

// heavily inspired by https://github.com/TTK4145/Project-resources/tree/master/cost_fns
type HRAElevState struct {
	Behavior    string           `json:"behaviour"`
	Floor       int              `json:"floor"`
	Direction   string           `json:"direction"`
	CabRequests [NUM_FLOORS]bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [NUM_FLOORS][NUM_HALL_BUTTONS]bool `json:"hallRequests"`
	States       map[string]HRAElevState            `json:"states"`
}

var behaviourToString = map[ElevatorBehaviour]string{
	EB_Idle:     "idle",
	EB_DoorOpen: "doorOpen",
	EB_Moving:   "moving"}

var directionToString = map[MotorDirection]string{
	MD_Up:   "up",
	MD_Down: "down",
	MD_Stop: "stop"}

var numberToString = map[int]string{
	1: "one",
	2: "two",
	3: "three",
	4: "four",
	5: "five",
	6: "six",
	7: "seven",
	8: "eight",
	9: "nine",
}

func CostFunction(
	hallOrders [NUM_FLOORS][NUM_HALL_BUTTONS]OrderType,
	cabOrders [NUM_ELEVATORS][NUM_FLOORS]bool,
	elevators [NUM_ELEVATORS]Elevator,
	desiredNodeID int,
	isAlive [NUM_ELEVATORS]bool,
) [NUM_FLOORS][NUM_BUTTONS]bool {
	defer phoenix.RespawnAfterPanic()

	var (
		assignedOrders [NUM_FLOORS][NUM_BUTTONS]bool
		hallOrdersBool [NUM_FLOORS][NUM_HALL_BUTTONS]bool
	)
	for floor := range assignedOrders {
		assignedOrders[floor] = [NUM_BUTTONS]bool{false, false, cabOrders[desiredNodeID][floor]}
	}

	var checkIfAnyHallOrders bool = false
	for floor := range hallOrders {
		for button, value := range hallOrders[floor] {
			hallOrdersBool[floor][button] = (value == ConfirmedOrder)
			if value == ConfirmedOrder {
				checkIfAnyHallOrders = true
			}
		}
	}

	if !checkIfAnyHallOrders {
		return assignedOrders
	}

	var aliveNodes int = 0
	for _, alive := range isAlive {
		if alive {
			aliveNodes++
		}
	}

	if aliveNodes == 0 {
		for floor := range assignedOrders {
			assignedOrders[floor][BT_HallUp] = hallOrdersBool[floor][BT_HallUp]
			assignedOrders[floor][BT_HallDown] = hallOrdersBool[floor][BT_HallDown]
		}
		return assignedOrders
	}

	input := HRAInput{
		HallRequests: hallOrdersBool,
		States:       make(map[string]HRAElevState, aliveNodes),
	}
	for id, elevator := range elevators {
		if !isAlive[id] {
			continue
		}
		input.States[numberToString[id+1]] = HRAElevState{
			Behavior:    behaviourToString[elevator.Behaviour],
			Floor:       elevator.Floor,
			Direction:   directionToString[elevator.Direction],
			CabRequests: cabOrders[id],
		}
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		panic("json.Marshal error: ")
	}

	_, path, _, _ := runtime.Caller(0)
	path = filepath.Dir(path)

	ret, err := exec.Command(path+COST_FUNCTION_COMPILED, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		panic(err)
	}

	output := make(map[string][][NUM_HALL_BUTTONS]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		panic("json.Unmarshal error: ")
	}

	var desiredNodeID_str = numberToString[desiredNodeID+1]
	for floor := range hallOrders {
		assignedOrders[floor][BT_HallUp] = output[desiredNodeID_str][floor][BT_HallUp]
		assignedOrders[floor][BT_HallDown] = output[desiredNodeID_str][floor][BT_HallDown]
	}

	return assignedOrders
}
