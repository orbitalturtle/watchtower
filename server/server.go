package main

import (
	"bufio"
        "context"
	"fmt"
	"net"
	"os"
	"strings"
        "sync"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type server struct {
	db    *db
        blockscanner *blockscanner

	peers map[*net.Addr]bool
}

func newServer(db *db) *server {
	return &server{
		db:    db,
		peers: make(map[*net.Addr]bool),
	}
}

func startServer(wg *sync.WaitGroup) {
	db, err := setUpDatabase()
	if err != nil {
		fmt.Println("Error setting up mongoDB: ", err)
	}
        defer db.client.Disconnect(context.TODO())

	s := newServer(db)

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

        wg.Done()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()

		// Handle connections in a new goroutine.
		go s.handleRequest(&conn)
	}
}

func (s *server) handleRequest(conn *net.Conn) {
	remoteAddr := (*conn).RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(*conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		s.handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func (s *server) handleMessage(cmd string, conn *net.Conn) {
	newCmd := strings.Trim(cmd, "\n ")

	if len(newCmd) > 0 {
		switch newCmd {
		case "/appointment":
			fmt.Println("Trying to add appointment to watchtower")
			s.addAppointment(conn)
		case "/init":
			fmt.Println("Initializing watchtower connection")
			s.initWatch(conn)
		default:
			(*conn).Write([]byte("Unrecognized command.\n"))
		}
	}
}

// TODO: Split off into a towerd file.
func main() {
	startServer(nil)
}
