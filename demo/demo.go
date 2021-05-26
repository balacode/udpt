// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                      /demo/[demo.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/balacode/udpt"
)

// main demo
func main() {
	// secret encryption key shared by the Sender and Receiver
	cryptoKey := []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	cf := udpt.NewDebugConfig() // verbose logging
	//
	// set-up and run the receiver
	const tag = "-------------> DEMO"
	fmt.Println(tag, "Running the receiver")
	received := ""
	rc := udpt.Receiver{
		Port:      9876,
		CryptoKey: cryptoKey,
		Config:    cf,
		//
		// receives fully-transferred data items sent to the receiver
		ReceiveData: func(k string, v []byte) error {
			received = string(v)
			div := strings.Repeat("##", 40)
			fmt.Println(tag, div)
			fmt.Println(tag, "You should see a 'Hello World!' message below:")
			fmt.Println(tag, div)
			fmt.Println(tag, "Receiver's ReceiveData received",
				"k:", k, "v:", received)
			fmt.Println(tag, div)
			return nil
		},
		// provides existing data items for hashing by the Receiver. Only the
		// hash will be sent back to the sender, to confirm the transfer.
		ProvideData: func(k string) ([]byte, error) {
			fmt.Println(tag, "Receiver's ProvideData()")
			return []byte(received), nil
		},
	}
	go func() { _ = rc.Run() }()
	//
	// send a message to the receiver
	time.Sleep(time.Second)
	fmt.Println(tag, "Sending a message")
	sender := udpt.Sender{
		Address: "127.0.0.1:9876", CryptoKey: cryptoKey, Config: cf,
	}
	err := sender.SendString("demo_data", "Hello World!")
	if err != nil {
		fmt.Println(tag, "failed sending:", err)
	}
	time.Sleep(time.Second)
} //                                                                        main

// end
