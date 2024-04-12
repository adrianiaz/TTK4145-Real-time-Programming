package elevio

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const _pollRate = 20 * time.Millisecond

var _initialized bool = false
var _numFloors int = 4
var _mtx sync.Mutex
var _conn net.Conn

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

// sets up a TCP connection to the elevator server
// initializes the elevator hardware
func Init(addr string, numFloors int) {
	if _initialized {
		fmt.Println("Driver already initialized!")
		return
	}
	_numFloors = numFloors
	_mtx = sync.Mutex{}
	var err error
	_conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	_initialized = true
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

// check if the button is pressed and send it to the channel
func PollButtons(receiver chan<- ButtonEvent) {
	prev := make([][3]bool, _numFloors) // [floor][button]
	for {
		time.Sleep(_pollRate)
		for f := 0; f < _numFloors; f++ { //for each floor
			for b := ButtonType(0); b < 3; b++ { //start at 0, end at 3, buttonType(0)= HallUP, buttonType(1)= HallDown, buttonType(2)= Cab
				v := GetButton(b, f) //check if button is pressed - check every button
				if v != prev[f][b] && v != false {
					receiver <- ButtonEvent{f, ButtonType(b)} //send to channel
				}
				prev[f][b] = v //update prev
			}
		}
	}
}

// check if the floor sensor is pressed and send it to the channel
// checks which floor the elevator is at, unless it is between floors and sends it wheen it reaches a floor
func PollFloorSensor(receiver chan<- int) {
	prev := -1 //
	for {
		time.Sleep(_pollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			receiver <- v
		}
		prev = v
	}
}

// check if the stop button is pressed and send it to the channel
func PollStopButton(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := GetStop()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

// check if the obstruction switch is pressed and send it to the channel
func PollObstructionSwitch(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := GetObstruction()
		if v != prev {
			receiver <- v
		}
		prev = v
	}
}

// return the value of the button, true or false
func GetButton(button ButtonType, floor int) bool {
	a := read([4]byte{6, byte(button), byte(floor), 0})
	return toBool(a[1])
}

// return the value of the floor sensor, true or false - where the elevator is
// if the elevator is between floors, it returns -1
// use this if we need to know which floor the elevator is at
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

// return the value of the obstruction switch, true or false
func GetObstruction() bool {
	a := read([4]byte{9, 0, 0, 0})
	return toBool(a[1])
}

// help function to read and write to the server
// sending a request to the server and getting a response
func read(in [4]byte) [4]byte {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	var out [4]byte
	_, err = _conn.Read(out[:])
	if err != nil {
		panic("Lost connection to Elevator Server")
	}

	return out
}

// help function to read and write to the server
// sending a command to the elevator server
func write(in [4]byte) {
	_mtx.Lock()
	defer _mtx.Unlock()

	_, err := _conn.Write(in[:])
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


func DirToString(dir MotorDirection) string {
	switch dir {
	case MD_Up:
		return "up"
	case MD_Down:
		return "down"
	case MD_Stop:
		return "stop"
	}
	return "unknown"
}

func ButtonToString(button ButtonType) string {
	switch button {
	case BT_HallUp:
		return "hall_up"
	case BT_HallDown:
		return "hall_down"
	case BT_Cab:
		return "cab"
	}
	return "unknown"
}