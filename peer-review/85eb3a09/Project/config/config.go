package config

const NumFloors = 4
const NumButtons = 3
const DoorOpenDuration = 3
const MotorStopDuration = 4
const StateUpdatePeriodMs = 500
const SendMessageTimerMs = 20

type RequestState int

const (
	None      RequestState = 0
	Order     RequestState = 1
	Confirmed RequestState = 2
	Complete  RequestState = 3
)

const (
	P_PEERS = 15669
	P_BCAST = 16569
)