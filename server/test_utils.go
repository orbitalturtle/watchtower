package main

import (
        "context"
        "testing"
)

func setUpTestDatabase(t *testing.T) (*db, func(), error) {
        testDb, closeDb, err := setUpDatabase() 
        if err != nil {
                t.Fatal("Encountered an error creating a mongo database")
                return nil, nil, err
        }

        deleteDb := func() {
            closeDb()
            testDb.db.Drop(context.TODO())
        }

        return testDb, deleteDb, nil 
}

