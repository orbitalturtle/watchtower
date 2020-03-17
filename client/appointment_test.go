package client

import (
        "reflect"
        "testing"
)

var ( 
        testAppointment = appointment{
                // TODO: Replace with an actual locator
                Locator:          "LocatorTest",
                StartBlock:       606680,
                EndBlock:         606690,
        } 

        testBlob = "1971df546d4cd81724dce4082c85a934271761268c3dfa3d197e2a8425bf4f73"

        testKey = "3ac41a3ab9dd598f61104ca4cbf840ec" 

        testData = secretBackupData{blob: []byte(testBlob)}
)

// Test that the appointment endpoint responds correctly.
func TestEncryptAppointment(t *testing.T) {
        _, err := testAppointment.encrypt(testKey, &testData)
        if err != nil {
                t.Fatal("Encrypting appointment produced err: ", err)
        }

        // TODO: Test that if the method is given bad data, it responds appropriately.
}

func TestAddEncryptedBlob(t *testing.T) {
        err := testAppointment.addEncryptedBlob(testKey, &testData)
        if err != nil {
                t.Fatal("Error adding encrypted blob to appointment: ", err)
        }

        xType := reflect.TypeOf(testAppointment.EncryptedBlob)

        if xType.String() != "[]uint8" {
                 t.Fatal("Encrypted blob was not set appropriately.")
        }
}
