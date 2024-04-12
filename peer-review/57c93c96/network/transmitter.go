package network

import (
	. "Project/config"
	"Project/phoenix"
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

func transmitter(payloadFromParentChannel <-chan PayloadOnNetwork) {
	defer phoenix.RespawnAfterPanic()
	conn, err := net.ListenPacket("udp4", "")
	if err != nil {
		log.Print(err)
		panic(err)
	}

	address, err := net.ResolveUDPAddr("udp4", NETWORK_BCAST_ADDR+NETWORK_PORT)
	if err != nil {
		panic(err)
	}

	for {
		message := encode_payload(<-payloadFromParentChannel)
		_, _ = conn.WriteTo(message, address)

		time.Sleep(NETWORK_BCAST_INTERVAL)
	}
}

func encode_payload(p PayloadOnNetwork) []byte {
	defer phoenix.RespawnAfterPanic()
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
