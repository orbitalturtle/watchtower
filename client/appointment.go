package client

import (
        "crypto/cipher"
        cryptorand "crypto/rand"

        "golang.org/x/crypto/chacha20poly1305"
)

// TODO: Should be more aligned with object-oriented principles.

type secretBackupData struct {
        chanID string
}

type origAppointment struct {
        Locator          string 
        StartBlock       uint64
        EndBlock         uint64
}

type appointment struct {
        Locator          string 
        StartBlock       uint64
        EndBlock         uint64
        EncryptedBlob    []byte
}

func buildAppointment(origAppt *originalAppointment) *appointment {
        // Take secret data and encrypt it

        // Add the encrypted blob to the appointment

        return appt
}

func encryptAppointment(key string, data *secretData) []byte, error {
        // Key is the second half of the tx data...
        
        // Convert data into a byte array so it can be encrypted.

        aeadCipher, err := chacha20poly1305.New([]byte(key))
        if err != nil {
               return nil, err
        }

        // Encryption
        nonce := make([]byte, chacha20poly1305.NonceSize)
        if _, err := cryptorand.Read(nonce); err != nil {
            panic(err)
        }
        encryptedBlob := aead.Seal(nil, nonce, []byte(data), nil)

        return encryptedBlob, nil
}
