package elevatordriver

import (
	. "Project/config"
	. "Project/datatypes"
)

// heavily inspired by https://github.com/TTK4145/Project-resources/blob/master/elev_algo/requests.c
func anyRequestsAbove(c Controller) bool {
	for floor := c.elevator.Floor + 1; floor < NUM_FLOORS; floor++ {
		for _, value := range c.orders[floor] {
			if value {
				return true
			}
		}
	}
	return false
}

func anyRequestsBelow(c Controller) bool {
	for floor := 0; floor < c.elevator.Floor; floor++ {
		for _, value := range c.orders[floor] {
			if value {
				return true
			}
		}
	}
	return false
}

func anyRequestsHere(c Controller) bool {
	for _, value := range c.orders[c.elevator.Floor] {
		if value {
			return true
		}
	}
	return false
}

func hasHallDownOrCabRequest(c Controller) bool {
	return c.orders[c.elevator.Floor][BT_HallDown] ||
		c.orders[c.elevator.Floor][BT_Cab]
}

func hasHallUpOrCabRequest(c Controller) bool {
	return c.orders[c.elevator.Floor][BT_HallUp] ||
		c.orders[c.elevator.Floor][BT_Cab]
}
