package main

import (
	"bufio"
        "context"
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

	var cmd string
	var response interface{}

	// TODO: Check that authentication is valid.
	// Make sure appointment message has every required property.
	if appointment.Locator == nil {
		// TODO: Check that locator is the correct length
		cmd = "AppointmentRejected \n"
		reason := "Locator is required. Must be the correct size to match a real transaction."
		response = AppointmentRejected{
			Locator:   appointment.Locator,
			Rcode:     400,
			Reason:    reason,
			ReasonLen: uint16(len(reason)),
		}
	} else {
		// If we hit none of the above errors, appointment message is acceptable.
		// Add peer to watchtower's list.
		peerAddr := (*conn).RemoteAddr()
		s.peers[&peerAddr] = true

                // Add appointment to our database.
                collection := s.db.Database("test").Collection("appointments")
                insertResult, err := collection.InsertOne(context.TODO(), appointment)
                if err != nil {
                    log.Fatal(err)
                }
                
                log.Println("Inserted an appointment into db: ", insertResult.InsertedID)

		cmd = "AppointmentAccepted \n"
		// Then send back an "AppointmentAccepted" response.
		response = AppointmentAccepted{
			Locator: appointment.Locator,
			Qos:     appointment.QosData,
			QosLen:  appointment.QosLen,
		}
	}

	rw.Writer.WriteString(cmd)
	rw.Flush()

	// Send an encoded response to the client.
	enc := gob.NewEncoder(rw)
	err = enc.Encode(response)
	if err != nil {
		log.Println("Error encoding GOB data:", err)
		return
	}
	rw.Flush()
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
