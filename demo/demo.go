// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /demo/[demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"fmt"
	"time"

	"github.com/balacode/udpt"
)

// main demo
func main() {
	// secret encryption key shared by the Sender and Receiver
	cryptoKey := []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// set-up and run the receiver
	received := ""
	rc := udpt.Receiver{Port: 9876, CryptoKey: cryptoKey,
		//
		// receives fully-transferred data items sent to the receiver
		ReceiveData: func(k string, v []byte) error {
			fmt.Println("Receiver.ReceiveData k:", k, "v:", string(v))
			received = string(v)
			return nil
		},
	}
	go func() { _ = rc.Run() }()
	defer func() { rc.Stop() }()
	time.Sleep(500 * time.Millisecond)
	//
	// send a message to the receiver
	err := udpt.SendString("127.0.0.1:9876", "msg", "Hello World!", cryptoKey)
	if err != nil {
		fmt.Println("failed sending:", err)
	}
	time.Sleep(500 * time.Millisecond)
} //                                                                        main

// end
