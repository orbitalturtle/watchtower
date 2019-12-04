package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

// Endpoint for scheduling a watchtower appointment.
func (s *server) addAppointment(conn *net.Conn) {
	var appointment Wt_appointment

	rw := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))

	dec := gob.NewDecoder(rw)
	err := dec.Decode(&appointment)
	if err != nil {
		log.Println("Error decoding appointment GOB data:", err)
		return
	}

	// TODO: Make sure appointment message has every required property.

	// If init message is acceptable, add peer to watchtower's list
	peerAddr := (*conn).RemoteAddr()
	s.peers[&peerAddr] = true

	var response bytes.Buffer
	// Send an encoded response to the client.
	enc := gob.NewEncoder(rw)
	err = enc.Encode(&response)
	if err != nil {
		log.Println("Error encoding GOB data:", err)
		return
	}
}

// Endpoint for init message, which kicks off communication with the watchtower.
func (s *server) initWatch(conn *net.Conn) {
	var init Wt_init

	rw := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))

	dec := gob.NewDecoder(rw)
	err := dec.Decode(&init)
	if err != nil {
		log.Println("Error decoding init GOB data:", err)
		return
	}
}
