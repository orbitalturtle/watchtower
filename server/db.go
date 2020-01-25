package main 

import (
    "context"
    "fmt"
    "log"

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

type db struct {
    client *mongo.Client
    db *mongo.Database
}

func setUpDatabase() (*db, func(), error) {
    // Set client options
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        fmt.Println("mongo connect error: ", err)
        return nil, nil, err
    }

    db := db{
        client: client,
        db: client.Database("test"),
    }

    fmt.Println("Connected to MongoDB!")

    err = db.createApptCollection() 
    if err != nil {
        fmt.Println("Create appt collection err: ", err)
        return nil, nil, err
    }

    closeDb := func() {
        db.client.Disconnect(context.TODO())
    }

    return &db, closeDb, nil
}

// Create collection with an index so watchtower can query for transactions faster.
func (d *db) createApptCollection() error {
    apptsCollection := d.client.Database("test").Collection("appointments")

    indexView := apptsCollection.Indexes()

    // Look through indexes to see if Locator_1 index already exists.
    // If it does, we don't need to create it again.
    cursor, err := indexView.List(context.TODO())
    if err != nil {
        log.Fatal("Problem listening appointment collection indexs: ", err)
      
    }

    // Get a slice of all indexes returned 
    var results []bson.M
    if err = cursor.All(context.TODO(), &results); err != nil {
        log.Fatal(err)
    }

    for _, index := range results {
        if index["name"] == "Locator_1" {
            return nil
        }
    }

    model := mongo.IndexModel{Keys: bson.D{{"Locator", 1}}}

    names, err := indexView.CreateOne(context.TODO(), model)
    if err != nil {
        fmt.Println("Err with creating appointments index: ", err)
        return err
    }

    fmt.Printf("Created appointment index %v\n", names)
    return nil
}

// Insert an appointment.
func (d *db) insertAppt(appointment Wt_appointment) error {
    collection := d.db.Collection("appointments")
    insertResult, err := collection.InsertOne(context.TODO(), appointment)
    if err != nil {
        fmt.Println("Error inserting appointment into db: ", err)
        return err
    }

    fmt.Println("Inserted an appointment into db: ", insertResult.InsertedID)
    return nil
}

func (d *db) deleteAppt(locator string) error {
    collection := d.db.Collection("appointments")
    _, err := collection.DeleteOne(context.TODO(), bson.D{{"locator", locator}})
    if err != nil {
        fmt.Println("Error deleting appointment from db: ", err)
        return err
    }

    fmt.Println("Deleted an appointment from db: ", locator)
    return nil
}
