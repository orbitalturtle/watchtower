package main

import (
    "testing"
)

func TestGetLocatorFromID(t *testing.T) {
    testTxid := "1325c73b1f8d0a6075cd063c98a528b04bbc01ec4d22695295a674f2caeace14"
    testLocator := "1325c73b1f8d0a6075cd063c98a528b0"

    locator := getLocatorFromTxid(testTxid)
    if locator != testLocator {
        t.Fatal("getLocatorFromTxId should have returned: ", testLocator)
    } 
}

