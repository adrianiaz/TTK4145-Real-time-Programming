package config

import (
	"time"
)

const (
	NUM_FLOORS       int = 4
	NUM_BUTTONS      int = 3
	NUM_HALL_BUTTONS int = 2
	NUM_ELEVATORS    int = 3

	COST_FUNCTION_COMPILED string = "/cost_function_assigner"
	NETWORK_BCAST_ADDR     string = "10.100.23.255"
	NETWORK_PORT           string = ":20027"
	ELEVATOR_PORT          string = "localhost:15657"

	POLL_RATE_DRIVER       time.Duration = 20 * time.Millisecond
	NETWORK_BCAST_INTERVAL time.Duration = 20 * time.Millisecond
	NETWORK_TIMEOUT        time.Duration = 1000 * time.Millisecond

	TIMEOUT_AFTER_RESPAWN       time.Duration = 10 * time.Second
	ELEVATOR_DOOR_OPEN_TIME     time.Duration = 3 * time.Second
	DECLARE_THREAD_DEAD_AFTER   time.Duration = 10 * time.Second
	ALLOWED_TIME_DOOR_OPEN      time.Duration = 10 * time.Second
	ALLOWED_TIME_BETWEEN_FLOORS time.Duration = 4 * time.Second
)
