package obstruction

import (
	"elevatorinterface"
	"time"
)

func Poll(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(elevatorinterface.PollRate)
		v := getObstruction()
		if v {
		}
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func getObstruction() bool {
	a := elevatorinterface.Read([4]byte{9, 0, 0, 0})
	return elevatorinterface.ToBool(a[1])
}
