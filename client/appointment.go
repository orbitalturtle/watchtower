package client

import (
        // "crypto/cipher"
        cryptorand "crypto/rand"
        "encoding/json"

        "golang.org/x/crypto/chacha20poly1305"
)

// TODO: Figure out what other data is needed here.
type secretBackupData struct {
        blob []byte
}

type appointment struct {
        Locator          string 
        StartBlock       uint64
        EndBlock         uint64
        EncryptedBlob    []byte
}

func (a *appointment) addEncryptedBlob(key string, blobData *secretBackupData) error {
        // Encrypt secret data.
        encryptedBlob, err := a.encrypt(key, blobData)
        if err != nil {
                return err
        }

        // Add the encrypted blob to the appointment.
        a.EncryptedBlob = encryptedBlob

        return nil
}

// Encrypt appointment before sending it to the watchtower so the watchtower 
// can't see all of our transaction data.
// Key sis the second half of the tx id
func (a *appointment) encrypt(key string, blobData *secretBackupData) ([]byte, error) {
        // Convert data into a byte array so it can be encrypted.
        apptBytes, _ := json.Marshal(*blobData)

        aeadCipher, err := chacha20poly1305.New([]byte(key))
        if err != nil {
               return nil, err
        }

        // Encrypt data so that only the second half of the tx data can open it.
        nonce := make([]byte, chacha20poly1305.NonceSize)
        if _, err := cryptorand.Read(nonce); err != nil {
            panic(err)
        }
        encryptedBlob := aeadCipher.Seal(nil, nonce, apptBytes, nil)

        return encryptedBlob, nil
}
