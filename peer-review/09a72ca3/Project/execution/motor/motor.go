// Adapted from https://github.com/TTK4145/driver-go

package motor

import (
	"elevatorglobals"
	"elevatorinterface"
)

func SetDirection(dir elevatorglobals.Direction) {
	elevatorinterface.Write([4]byte{1, byte(dir), 0, 0})
}
