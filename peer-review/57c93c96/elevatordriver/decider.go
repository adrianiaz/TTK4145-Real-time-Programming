package elevatordriver

import (
	. "Project/datatypes"
)

func updateDirectionAndBehaviour(controller Controller) Elevator {
	var todo DirectionBehaviourPair = chooseDirection(controller)
	return Elevator{
		Floor:     controller.elevator.Floor,
		Direction: todo.direction,
		Behaviour: todo.behavior,
	}
}

// heavily inspired by https://github.com/TTK4145/Project-resources/blob/master/elev_algo/requests.c
func chooseDirection(c Controller) DirectionBehaviourPair {

	switch c.elevator.Direction {

	case MD_Up:
		if anyRequestsAbove(c) {
			return DirectionBehaviourPair{MD_Up, EB_Moving}
		} else if anyRequestsHere(c) {
			return DirectionBehaviourPair{MD_Down, EB_DoorOpen}
		} else if anyRequestsBelow(c) {
			return DirectionBehaviourPair{MD_Down, EB_Moving}
		}
		return DirectionBehaviourPair{MD_Stop, EB_Idle}

	case MD_Down:
		if anyRequestsBelow(c) {
			return DirectionBehaviourPair{MD_Down, EB_Moving}
		} else if anyRequestsHere(c) {
			return DirectionBehaviourPair{MD_Up, EB_DoorOpen}
		} else if anyRequestsAbove(c) {
			return DirectionBehaviourPair{MD_Up, EB_Moving}
		}
		return DirectionBehaviourPair{MD_Stop, EB_Idle}

	case MD_Stop:
		if anyRequestsHere(c) {
			return DirectionBehaviourPair{MD_Stop, EB_DoorOpen}
		} else if anyRequestsAbove(c) {
			return DirectionBehaviourPair{MD_Up, EB_Moving}
		} else if anyRequestsBelow(c) {
			return DirectionBehaviourPair{MD_Down, EB_Moving}
		}
		return DirectionBehaviourPair{MD_Stop, EB_Idle}

	default:
		return DirectionBehaviourPair{MD_Stop, EB_Idle}
	}
}

func shouldStop(c Controller) bool {
	switch c.elevator.Direction {
	case MD_Down:
		return hasHallDownOrCabRequest(c) || !anyRequestsBelow(c)

	case MD_Up:
		return hasHallUpOrCabRequest(c) || !anyRequestsAbove(c)

	default:
		return true
	}
}
