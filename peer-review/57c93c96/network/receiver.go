package network

import (
	. "Project/config"
	"Project/phoenix"
	"bytes"
	"encoding/gob"
	"net"
	"reflect"
	"time"
)

func receiver(
	payloadToParent chan<- PayloadOnNetwork,
	connectedListToParent chan<- [NUM_ELEVATORS]bool,
	myID int,
) {
	defer phoenix.RespawnAfterPanic()
	var (
		message          PayloadOnNetwork
		recentPayloads   [NUM_ELEVATORS]PayloadOnNetwork
		lastHeardFrom    [NUM_ELEVATORS]time.Time
		isConnectedList  [NUM_ELEVATORS]bool
		wasConnectedList [NUM_ELEVATORS]bool
		buffer           []byte = make([]byte, 1024)
	)

	conn, err := net.ListenPacket("udp4", NETWORK_PORT)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		wasConnectedList = isConnectedList

		conn.SetReadDeadline(time.Now().Add(NETWORK_BCAST_INTERVAL))
		_, _, err := conn.ReadFrom(buffer)
		if err == nil {
			message = decodePayload(buffer)
			if message.Id != myID {
				lastHeardFrom[message.Id] = time.Now()
				lastHeardFrom[myID] = time.Now()
			}
			if !reflect.DeepEqual(message, recentPayloads[message.Id]) {
				recentPayloads[message.Id] = message
				payloadToParent <- message
			}
		}

		isConnectedList = [NUM_ELEVATORS]bool{}
		for nodeID, lastLifeSignal := range lastHeardFrom {
			if lastLifeSignal.Add(NETWORK_TIMEOUT).After(time.Now()) {
				isConnectedList[nodeID] = true
			}
		}

		if !isConnectedList[myID] {
			recentPayloads = [NUM_ELEVATORS]PayloadOnNetwork{}
		}

		if !reflect.DeepEqual(isConnectedList, wasConnectedList) {
			connectedListToParent <- isConnectedList
		}
	}
}

func decodePayload(msg []byte) PayloadOnNetwork {
	defer phoenix.RespawnAfterPanic()
	payload := PayloadOnNetwork{}
	dec := gob.NewDecoder(bytes.NewReader(msg))
	err := dec.Decode(&payload)
	if err != nil {
		panic(err)
	}
	return payload
}
