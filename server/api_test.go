package main

import (
	"bufio"
	"encoding/gob"
	"net"
	"testing"
	"time"
)

func initServer() {
	main()
}

// Test that the init endpoint responds if correct data is passed to it.
func TestAppointmentEndpoint(t *testing.T) {
	go initServer()

	// Give server time to get set up
	time.Sleep(time.Second)

	cmd := "/appointment"

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

	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		t.Fatalf("Some problem connecting to watchtower: %v", err)
	}
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	rw.Writer.Write([]byte(cmd))

	// var initBytes bytes.Buffer
	enc := gob.NewEncoder(rw)
	err = enc.Encode(appointment)
	if err != nil {
		t.Fatalf("Encoding error: %v", err)
	}
}

// Test that the init endpoint responds if correct data is passed to it.
// func TestInitEndpoint(t *testing.T) {
// 	go initServer()
//
// 	// Give server time to get set up
// 	time.Sleep(time.Second)
//
// 	cmd := "/init"
//
// 	init := Wt_init{
// 		AcceptedCiphers: []string{"chacha20"},
// 		Modes:           []Mode{Altruistic},
// 		Qos:             []string{"accountability"},
// 	}
//
// 	conn, err := net.Dial("tcp", "127.0.0.1:3333")
// 	if err != nil {
// 		t.Fatalf("Some problem connecting to watchtower: %v", err)
// 	}
// 	defer conn.Close()
//
// 	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
//
// 	rw.Writer.Write([]byte(cmd))
//
// 	// var initBytes bytes.Buffer
// 	enc := gob.NewEncoder(rw)
// 	err = enc.Encode(init)
// 	if err != nil {
// 		t.Fatalf("Encoding error: %v", err)
// 	}
// }
