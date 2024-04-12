package motor

import (
	"elevatorglobals"
	"time"
)

func RunWatchdog(watchdogTimeoutChannel chan<- bool, directionUpdateChannel <-chan elevatorglobals.Direction, watchdogFloorUpdateChannel <-chan bool, watchdogObstructionRemovedChannel <-chan bool) {
	motorWatchdogDuration := 7 * time.Second
	watchdogTicker := time.NewTicker(motorWatchdogDuration)
	for {
		select {
		case <-watchdogTicker.C:
			watchdogTimeoutChannel <- true			
		case direction := <-directionUpdateChannel:
			if direction == elevatorglobals.Direction_Stop {
				watchdogTicker.Stop()
			} else {
				watchdogTicker = time.NewTicker(motorWatchdogDuration)
			}
		case <-watchdogFloorUpdateChannel:
			watchdogTicker = time.NewTicker(motorWatchdogDuration)
		case <-watchdogObstructionRemovedChannel:
			watchdogTicker = time.NewTicker(motorWatchdogDuration)
		}
	}					
}