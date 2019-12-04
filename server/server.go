package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

// TODO: Split off into a towerd file.
func main() {
	db, err := setUpDatabase()
	if err != nil {
		fmt.Println("Error setting up mongoDB: ", err)
	}

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

type server struct {
	db    *mongo.Client
	peers map[*net.Addr]bool
}

func newServer(db *mongo.Client) *server {
	return &server{
		db:    db,
		peers: make(map[*net.Addr]bool),
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
	if len(cmd) > 0 {
		switch cmd {
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
