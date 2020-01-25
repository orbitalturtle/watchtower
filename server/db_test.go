package main

import (
        "context"
        "reflect"
        "testing"

        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
)

func TestSetUpDatabase(t *testing.T) {
        clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
        client, err := mongo.Connect(context.TODO(), clientOptions)

        // If "test" database exists, then delete it.
        result, err := client.ListDatabaseNames(context.TODO(), bson.D{})
        if err != nil {
                t.Fatal("List database names error: ", err)
        }
        
        for _, db := range result {
                if db == "test" {
                        err = client.Database("test").Drop(context.TODO())
                        if err != nil {
                                t.Fatal("Failed to drop test database")
                        }
                        t.Log("Old test database deleted")
                } 
        }

        // Test that creating the database from scratch works.
        testDb, deleteDb, err := setUpTestDatabase(t) 
        if err != nil {
                t.Fatal("Encountered an error creating a mongo database")
        }
        defer deleteDb()

        dbType := reflect.TypeOf(testDb).String()
        if dbType != "*main.db" {
                 t.Fatal("Db is of the wrong type")
        }

        testLocator := "0cf6a4e9dfa5456542b7775155f7162d"

        // Add one appointment to the db.
        appointment := Wt_appointment{Locator: testLocator}

        testDb.insertAppt(appointment)
        if err != nil {
                t.Fatal("Error inserting appointment")
        }

        // Test that if test database already exists, setUpTestDatabase doesn't
        // override it. 
        testDb, deleteDb, err = setUpTestDatabase(t) 
        if err != nil {
                t.Fatal("Encountered an error seting up database")
        }
        defer deleteDb()

        // Make sure that the appointment that was added has persisted now
        // that database has been restarted.
        var checkAppointment Wt_appointment

        err = testDb.db.Collection("appointments").FindOne(context.TODO(), bson.D{{"locator", testLocator}}).Decode(&checkAppointment)
        if err != nil {
                t.Fatal("Encountered an error finding appointment by locator")
        }
}

func TestInsertAppt(t *testing.T) {
        testDb, deleteDb, err := setUpTestDatabase(t) 
        if err != nil {
                t.Fatal("Encountered an error creating a mongo database")
        }
        defer deleteDb()

        testLocator := "7df6a4e9dfa5456542b7775155f7162d"

        // Add one appointment to the db.
        appointment := Wt_appointment{Locator: testLocator}

        testDb.insertAppt(appointment)
        if err != nil {
                t.Fatal("Error inserting appointment")
        }

        // Make sure that the appointment was really added to the database.
        var checkAppointment Wt_appointment

        err = testDb.db.Collection("appointments").FindOne(context.TODO(), bson.D{{"locator", testLocator}}).Decode(&checkAppointment)
        if err != nil {
                t.Fatal("Encountered an error finding appointment by locator")
        }
}

func TestDeleteAppt(t *testing.T) {
        testDb, deleteDb, err := setUpTestDatabase(t) 
        if err != nil {
                t.Fatal("Encountered an error creating a mongo database")
        }
        defer deleteDb()

        testLocator := "fe4a8c538d87cd8bd905ed2fbfc372a8"

        // Add one appointment to the db.
        appointment := Wt_appointment{Locator: testLocator}

        err = testDb.insertAppt(appointment)
        if err != nil {
                t.Fatal("Error inserting appointment")
        }

        err = testDb.deleteAppt(testLocator)
        if err != nil {
                t.Fatal("Error deleting appointment")
        }

        // Make sure that the appointment was really deleted from the database.
        var checkAppointment Wt_appointment

        err = testDb.db.Collection("appointments").FindOne(context.TODO(), bson.D{{"locator", testLocator}}).Decode(&checkAppointment)
        if err == nil {
                t.Fatal("FindOne should not have found a document. It should have been deleted")
        }
}
