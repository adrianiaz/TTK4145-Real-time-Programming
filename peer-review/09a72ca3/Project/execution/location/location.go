// Adapted from https://github.com/TTK4145/driver-go

package location

import (
	"elevatorinterface"
	"time"
)

func SetFloorIndicator(floor int) {
	elevatorinterface.Write([4]byte{3, byte(floor), 0, 0})
}

func Poll(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(elevatorinterface.PollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}

func GetFloor() int {
	a := elevatorinterface.Read([4]byte{7, 0, 0, 0})
	if a[1] != 0 {
		return int(a[2])
	} else {
		return -1
	}
}
