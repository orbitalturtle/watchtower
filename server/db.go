package main 

import (
    "context"
    "fmt"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Data about the latest commitment transaction for a particular client, which 
// the watchtower needs to send a justice transaction, if needed.
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

    err = createApptCollection(client) 
    if err != nil {
        fmt.Println("Create appt collection err: ", err)
    }

    return client, nil
}

// Create collection with an index so watchtower can query for transactions faster.
func createApptCollection(client *mongo.Client) error {
    apptsCollection := client.Database("test").Collection("appointments")

    indexView := apptsCollection.Indexes()

    model := mongo.IndexModel{Keys: bson.D{{"locator", 1}}}

    names, err := indexView.CreateOne(context.TODO(), model)
    if err != nil {
        fmt.Println("Err with creating appointments index: ", err)
        return err
    }

    fmt.Printf("Created appointment index %v\n", names)
    return nil
}
