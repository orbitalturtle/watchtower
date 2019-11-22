package main 

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "sync"
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
        defer conn.Close()

        // Handle connections in a new goroutine.
        go handleRequest(&conn)
    }
}

func handleRequest(conn *net.Conn) {
    remoteAddr := (*conn).RemoteAddr().String()
    fmt.Println("Client connected from " + remoteAddr)

    scanner := bufio.NewScanner(*conn)

    for {
    	ok := scanner.Scan()

    	if !ok {
    		break
    	}

    	handleMessage(scanner.Text(), conn)
    }

    fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(cmd string, conn *net.Conn) {
    if len(cmd) > 0 {
        switch { 
        case cmd == "/init":
            fmt.Println("Initializing watchtower connection")
            initWatch(conn)

        default:
            (*conn).Write([]byte("Unrecognized command.\n"))
        }
    }
}
