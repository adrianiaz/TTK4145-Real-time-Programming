package elevio

import (
	. "Project/config"
	. "Project/datatypes"
	"Project/phoenix"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	initialized bool = false
	conn        net.Conn
	mtx         sync.Mutex
)

func Init(addr string) {
	defer phoenix.RespawnAfterPanic()
	if initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	mtx = sync.Mutex{}
	var err error
	conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	initialized = true
}

func SetMotorDirection(dir MotorDirection) {
	write([4]byte{1, byte(dir), 0, 0})
}

func SetButtonLamp(button ButtonType, floor int, value bool) {
	write([4]byte{2, byte(button), byte(floor), toByte(value)})
}

func SetFloorIndicator(floor int) {
	write([4]byte{3, byte(floor), 0, 0})
}

func SetDoorOpenLamp(value bool) {
	write([4]byte{4, toByte(value), 0, 0})
}

func SetStopLamp(value bool) {
	write([4]byte{5, toByte(value), 0, 0})
}

func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][NUM_BUTTONS]bool, NUM_FLOORS)
	for {
		time.Sleep(POLL_RATE_DRIVER)
		for f := 0; f < NUM_FLOORS; f++ {
			for b := ButtonType(0); b < ButtonType(NUM_BUTTONS); b++ {
				v := GetButton(b, f)
				if v != prev[f][b] && v {
					receiver <- ButtonEvent{Floor: f, Button: ButtonType(b)}
				}
				prev[f][b] = v
			}
		}
	}
}

func PollFloorSensor(receiver chan<- int) {
	prev := -1
	for {
		time.Sleep(POLL_RATE_DRIVER)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}

func PollStopButton(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(POLL_RATE_DRIVER)
		v := GetStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

func PollObstructionSwitch(receiver chan<- bool, exit <-chan bool) {
	prev := false
	for {
		time.Sleep(POLL_RATE_DRIVER)
		v := GetObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
		select {
		case <-exit:
			return
		default:
		}
	}
}

func GetButton(button ButtonType, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

func GetFloor() int {
	a := read([4]byte{7, 0, 0, 0})
	if a[1] != 0 {
		return int(a[2])
	} else {
		return -1
	}
}

func GetStop() bool {
	a := read([4]byte{8, 0, 0, 0})
	return toBool(a[1])
}

func GetObstruction() bool {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
}

func read(in [4]byte) [4]byte {
	defer phoenix.RespawnAfterPanic()
	mtx.Lock()
	defer mtx.Unlock()

	_, err := conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	var out [4]byte
	_, err = conn.Read(out[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	return out
}

func write(in [4]byte) {
	defer phoenix.RespawnAfterPanic()
	mtx.Lock()
	defer mtx.Unlock()

	_, err := conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}
}

func toByte(a bool) byte {
	var b byte = 0
	if a {
		b = 1
	}
	return b
}

func toBool(a byte) bool {
	var b bool = false
	if a != 0 {
		b = true
	}
	return b
}
