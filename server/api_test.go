package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strings"
        "sync"
	"testing"
)

// Test that the appointment endpoint responds correctly.
func TestAppointmentEndpoint(t *testing.T) {
        var wg sync.WaitGroup
        wg.Add(1)

	go startServer(&wg)

        wg.Wait()
 
	cmd := "/appointment \n"

	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		t.Fatalf("Some problem connecting to watchtower: %v", err)
	}
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	appointmentErr := Wt_appointment{Locator: nil}

	_, err = rw.Writer.WriteString(cmd)
	if err != nil {
		t.Fatalf("Writing cmd throws an error: %v", err)
	}
	rw.Flush()

	enc := gob.NewEncoder(rw)
	err = enc.Encode(appointmentErr)
	if err != nil {
		t.Fatalf("Encoding error when sending appointmentErr: %v", err)
	}
	rw.Flush()

	cmd, err = rw.Reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Reading error when retrieving appointment response: %v", err)
	}

	cmd = strings.TrimSpace(cmd)

	if cmd != "AppointmentRejected" {
		t.Fatalf("Watchtower should have returned an appointment error.")
	}

	var rejected AppointmentRejected

	dec := gob.NewDecoder(rw)
	err = dec.Decode(&rejected)
	if err != nil {
		log.Println("Error decoding response GOB data:", err)
		return
	}

	rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	encryptedBlob := []byte("BlobTest")
	authToken := []byte("TokenTest")
	qos := "Accountability"

	appointment := Wt_appointment{
		// TODO: Replace this with actual locator information
		Locator:          []byte("LocatorTest"),
		StartBlock:       606680,
		EndBlock:         606690,
		EncryptedBlob:    encryptedBlob,
		EncryptedBlobLen: uint16(len(encryptedBlob)),
		AuthToken:        authToken,
		AuthTokenLen:     uint16(len(authToken)),
		QosData:          qos,
		QosLen:           uint16(len(qos)),
	}

	cmd = "/appointment \n"

	_, err = rw.Writer.WriteString(cmd)
	if err != nil {
		t.Fatalf("Writing cmd throws an error: %v", err)
	}
	rw.Flush()

	enc = gob.NewEncoder(rw)
	err = enc.Encode(appointment)
	if err != nil {
		t.Fatalf("Encoding error when sending appointment: %v", err)
	}
	rw.Flush()

	// Check that the response comes back if we send correct data.
	cmd, err = rw.Reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Reading error when retrieving appointment response: %v", err)
	}

	cmd = strings.TrimSpace(cmd)

	if cmd != "AppointmentAccepted" {
		t.Fatalf("Appointment server shouldn't have returned an error.")
	}
}
