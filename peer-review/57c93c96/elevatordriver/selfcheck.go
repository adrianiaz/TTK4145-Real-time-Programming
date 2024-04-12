package elevatordriver

import (
	. "Project/config"
	"time"
)

func selfCheck(controller Controller) bool {
	var status bool = true
	if controller.moving {
		if controller.lastKnownMovement.Add(ALLOWED_TIME_BETWEEN_FLOORS).Before(time.Now()) {
			status = false
		}
	} else if controller.doorOpen {
		if controller.mostRecentDoorOpening.Add(ALLOWED_TIME_DOOR_OPEN).Before(time.Now()) {
			status = false
		}
	}
	return status
}
