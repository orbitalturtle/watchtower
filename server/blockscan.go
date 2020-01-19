package main

import (
        "bytes"
        "context"
        "encoding/hex"
        "fmt"
        "log"

        "github.com/btcsuite/btcd/chaincfg/chainhash"
        "github.com/btcsuite/btcd/rpcclient"
        "github.com/btcsuite/btcd/wire"
        zmq "github.com/pebbe/zmq4"

        "go.mongodb.org/mongo-driver/bson"
        "go.mongodb.org/mongo-driver/mongo"
)

type blockscanner struct {
        subscriber *zmq.Socket

        db *db
}

type Block struct {
       hash string
       confirmations int32
       size int32
       strippedsize int32
       weight int32
       height int32
       version int32
       versionHex string
       merkle_root string
       tx []string
       time uint32
       mediumtime uint32
       nonce uint32
       bits string
       difficulty uint32
       chainwork string
       nTx int32 
       previousblockhash string
}

func (b *blockscanner) start(db *db) {
        b.db = db

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

// handleBlock is called when our node receives a new bitcoin block.
func (b *blockscanner) handleBlock() {
        fmt.Println("Collecting block updates from bitcoind")

        for {
        	blockBytes2D, err := b.subscriber.RecvMessageBytes(0)
        	if err != nil {
        	        fmt.Println("Receiving block err: ", err)
        	}

                // Turn 2D slice of bytes into a 1D slice
                blockHashBytesRough := bytes.Join(blockBytes2D, nil) 

                // Convert bytes to a hexadecimal string. 
                blockHashStr := hex.EncodeToString(blockHashBytesRough)

                // Cut away unnecessary data to get to blockHash
                blockHashStrStripped := blockHashStr[18:82]

                // Encode back to bytes. 
                blockHashBytes, err := hex.DecodeString(blockHashStrStripped)
                if err != nil {
                        log.Fatal("Err turning hex block hash into bytes")
                }

                blockHashBytesReversed := reverseByteSlice(blockHashBytes)
 
                blockHash, err := chainhash.NewHash(blockHashBytesReversed)
                if err != nil {
	                log.Fatal("Err generating block hash from bytes: ", err)
                }

                fmt.Println("blockHash= ", blockHash)

                client, err := connectToRPC()
                if err != nil {
                        log.Fatal(err)
                }
                defer client.Shutdown()

                block, err := client.GetBlock(blockHash)
	        if err != nil {
	        	log.Fatal("Err retrieving block info via RPC= ", err)
	        }

                fmt.Println("Received a new block from bitcoind: ", block)

                fraudTxs, err := b.lookForMatches(block.Transactions)
                if err != nil {
                        log.Fatal("Searching through block for matches err: ", err)
                        return 
                } 
                if fraudTxs == nil {
                        // There are no matches. Go back to waiting for a new block. 
                        return
                } else {
                        // Need to send a justice transaction.
                }
        }
}

// TODO: Should block arg here be a pointer?
// lookForMatches scans through all the new transaction ids in the new bitcoin block.
// It looks to see if any of them match the ids we have in our database. 
func (b *blockscanner) lookForMatches(txs []*wire.MsgTx) ([]string, error) {
        matches := make([]string, 0)

        // Loop through block transactions. 
        for _, tx := range txs {
            txid := tx.TxHash()

            // Grab the txid and turn it into a locator. 
            locator := getLocatorFromTxid(txid.String())
 
            // Then query locator index in database to see if it turns up a match.
            apptsCollection := b.db.client.Database("test").Collection("appointments")

            var appointment Wt_appointment

            err := apptsCollection.FindOne(context.TODO(), bson.D{{"locator", locator}}).Decode(&appointment)
            if err == nil {
                // Found a match
                log.Println("Found a match: ", err)
                matches = append(matches, locator)
            } else {
                // ErrNoDocuments means that the filter did not match any documents in the collection
                if err == mongo.ErrNoDocuments {
                    log.Println("No documents found wih this locator index: ", err)
                    continue 
                }
                log.Println("Error searching appointments collection via locator index: ", err)
                return nil, err
            }
        }

        if len(matches) == 0 {
            log.Println("No fraudulent transactions found in this block.")
            return nil, nil
        } else {
	    return matches, nil
        }
}

func reverseByteSlice(bytes []byte) []byte {
    for i, j := 0, len(bytes) - 1; i < j; i, j = i + 1, j - 1 {
        bytes[i], bytes[j] = bytes[j], bytes[i]
    }

    return bytes 
}

func connectToRPC() (*rpcclient.Client, error) {
        // Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "meh",
		Pass:         "meh",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
                return nil, err
	}

        

        return client, nil
}


