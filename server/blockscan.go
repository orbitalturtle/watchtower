package main

import (
        "fmt"
        "sync"

        zmq "github.com/pebbe/zmq4"
)

type blockscanner struct {
        subscriber *zmq.Socket
}

func (b *blockscanner) start() {
        //  Socket to receive block updates from bitcoind 
	subscriber, err := zmq.NewSocket(zmq.SUB)
        if err != nil {
                fmt.Println("ZMQ NewSocket err: ", err)
        }
        b.subscriber = subscriber
        defer subscriber.Close()

	err = b.subscriber.Connect("tcp://localhost:29000")
        if err != nil {
                fmt.Println("ZMQ connect err: ", err)
        }

        err = b.subscriber.SetSubscribe("hashblock")
        if err != nil {
		fmt.Println("ZMQ subscribe err: ", err)
        }
 
        b.handleBlock()
}

func (b *blockscanner) handleBlock() {
        fmt.Println("Collecting block updates from bitcoind")

        for {
        	newBlock, err := b.subscriber.RecvMessage(0)
        	if err != nil {
        	        fmt.Println("Receiving block err: ", err)
        	}
        }
}
