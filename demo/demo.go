// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /demo/[demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"fmt"

	"github.com/balacode/udpt"
)

// main demo
func main() {
	// secret encryption key shared by the Sender and Receiver
	cryptoKey := []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// set-up and run the receiver
	rc := udpt.Receiver{Port: 9876, CryptoKey: cryptoKey,
		Receive: func(k string, v []byte) error {
			fmt.Println("Received k:", k, "v:", string(v))
			return nil
		}}
	go func() { _ = rc.Run() }()
	//
	// send a message to the receiver
	err := udpt.SendString("127.0.0.1:9876", "main", "Hello World!", cryptoKey)
	if err != nil {
		fmt.Println("failed sending:", err)
	}
	rc.Stop()
} //                                                                        main

// end
