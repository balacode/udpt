// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /demo/[demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"context"
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
	defer demo2()
} //                                                                        main

// demo2 shows another way to run a Receiver, by using udpt.Receive function.
func demo2() {
	// secret encryption key shared by the Sender and Receiver
	cryptoKey := []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// run the receiver
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := udpt.Receive(ctx, 9876, cryptoKey,
			func(k string, v []byte) error {
				fmt.Println("Received k:", k, "v:", string(v))
				return nil
			})
		if err != nil {
			fmt.Println("failed receiving:", err)
		}
	}()
	// send a message to the receiver
	err := udpt.SendString("127.0.0.1:9876", "demo2", "Hello World!", cryptoKey)
	if err != nil {
		fmt.Println("failed sending:", err)
	}
	cancel()
} //                                                                       demo2

// end
