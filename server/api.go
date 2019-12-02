package main 

import (
    "bufio"
    "encoding/gob"
    "log"
    "net"
)

// Endpoint for init message, which kicks off communication with the watchtower. 
func (s *server) initWatch(conn *net.Conn) {
    var init Wt_init 

    rw := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))

    dec := gob.NewDecoder(rw)
    err := dec.Decode(&init)
    if err != nil {
	log.Println("Error decoding GOB data:", err)
        return
    }  
}
