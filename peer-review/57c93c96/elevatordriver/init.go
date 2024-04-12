package elevatordriver

import (
	. "Project/config"
	. "Project/datatypes"
)

func initializeController() Controller {
	return Controller{
		elevator: Elevator{
			Floor:     0,
			Direction: MD_Up,
			Behaviour: EB_Moving,
		},
		orders:      [NUM_FLOORS][NUM_BUTTONS]bool{},
		doorOpen:    false,
		moving:      false,
		sendStatus:  false,
		obstruction: false,
	}
}
