package main

import (
    "fmt"
    "net"
    "os"
) 

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {
    _, err := setUpDatabase()
    if err != nil {
        fmt.Println("Error setting up mongoDB: ", err)
    }   

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()

    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(&conn)
    }
}

func handleRequest(conn *net.Conn) {
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)

    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
      fmt.Println("Error reading:", err.Error())
    }

    // Convert initial bytes to string to determine which endpoint to use 
    json

    switch

    // Send a response back to person contacting us.
    conn.Write([]byte("Message received."))

    // Close the connection when done with it.
    conn.Close()

    fmt.Println("Server has received a connection: ", conn)
}
