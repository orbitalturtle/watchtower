package main

func getLocatorFromTxid(txid string) string {
        // Take first half of the txid to get the locator.
        // This is to get the first half (16 bytes) of the string. Do we want to keep it in hexadecimal?
        return txid[0:32]
}

