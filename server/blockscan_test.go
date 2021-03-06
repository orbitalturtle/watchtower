package main

import (
    "bytes"
    "fmt"
    "testing"
    "time"

    "github.com/btcsuite/btcd/wire"
)

var testTx = wire.MsgTx{
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash: [32]byte{0x01},
			},
		},
		{
			PreviousOutPoint: wire.OutPoint{
				Hash: [32]byte{0x00},
			},
		},
	},
}

var testTxs = []*wire.MsgTx{&testTx}

var testTxHash = "78d1fb1d07e48c8296c93995d8262bea1e5eb8d8bc1568df9fc5db82e676254d"

func startBlockscanner(t *testing.T) (*blockscanner, func()) {
	db, deleteDb, err := setUpTestDatabase(t)
        if err != nil {
                fmt.Println("Error setting up mongoDB: ", err)
        }

        s := newServer(db)

        blockscanner := &blockscanner{}
        s.blockscanner = blockscanner
        go (s.blockscanner).start(db)

        return s.blockscanner, deleteDb
}

func TestLookForMatches(t *testing.T) {
        b, deleteDb := startBlockscanner(t) 
        defer deleteDb()

        time.Sleep(1 * time.Second)

        appointment := Wt_appointment{
            Locator: getLocatorFromTxid(testTxHash),
        } 

        err := b.db.insertAppt(appointment)
        if err != nil {
            t.Fatal("Err inserting appointment: ", err)
        } 

        matches, err := b.lookForMatches(testTxs)
        if matches == nil {
            t.Fatal("This test tx should return as a match")
        }

        err = b.db.deleteAppt(appointment.Locator)
        if err != nil {
            t.Fatal(err)
        }

        // If locator is not in database, watchtower should find it is not a
        // match. 
        matches, err = b.lookForMatches(testTxs)
        if matches != nil {
            t.Fatal("Watchtower shouldn't find this tx as a match")
        }
}

func TestReverseByteSlice(t *testing.T) {
    origBytes := []byte{101, 8, 29, 90}
    expectedReversedBytes := []byte{90, 29, 8, 101}

    reversed := reverseByteSlice(origBytes)

    if !bytes.Equal(reversed, expectedReversedBytes) {
        t.Fatal("reverseByteSlice should have returned: ", expectedReversedBytes)
    }
}

