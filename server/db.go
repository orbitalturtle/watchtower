package main 

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Data about the latest commitment transaction for a particular client, which 
// the watchtower needs to send a justice transaction.
type appointment struct {
    locator []byte
    startBlock uint64
    endBlock uint64
    disputeDelta uint64
    encryptedBlob []byte
    transactionSize uint64
    transactionFee uint64
    cipher uint16
    op_customer_signature_algorithm uint16
    op_customer_signature []byte
    op_customer_public_key []byte
}

func setUpDatabase() (*mongo.Client, error) {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        fmt.Println("mongo connect error: ", err)
        return nil, err
    }

    fmt.Println("Connected to MongoDB!")

    return client, nil
}


