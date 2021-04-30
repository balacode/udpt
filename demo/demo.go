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

func main() {
	// the encryption key shared by the sender and receiver
	var cryptoKey = []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
	//
	// disable log buffering and enable verbose logging: for demos/debugging
	udpt.LogBufferSize = -1
	cfg := udpt.DefaultConfig()
	cfg.VerboseSender = true
	cfg.VerboseReceiver = true
	//
	// set-up and run the receiver
	prt("Running the receiver")
	receiver := udpt.Receiver{
		Port:      1234,
		CryptoKey: cryptoKey,
		Config:    cfg,
		//
		// receives fully-transferred data items sent to the receiver
		ReceiveData: func(name string, data []byte) error {
			div := strings.Repeat("##", 40)
			prt(div)
			prt("You should see a 'Hello World!' message below:")
			prt(div)
			prt("Receiver's write received name:", name, "data:", string(data))
			prt(div)
			return nil
		},
		// provides existing data items for hashing by the Receiver. Only the
		// hash will be sent back to the sender, to confirm the transfer.
		ProvideData: func(name string) ([]byte, error) {
			prt("Receiver's ProvideData()")
			return nil, nil
		},
	}
	go receiver.Run()
	//
	// send a message to the receiver
	time.Sleep(1 * time.Second)
	prt("Sending a message")
	sender := udpt.Sender{
		Address:   "127.0.0.1",
		Port:      1234,
		CryptoKey: cryptoKey,
		Config:    cfg,
	}
	err := sender.SendString("demo_data", "Hello World!")
	if err != nil {
		prt("Failed sending:", err)
	}
	wait := 30 * time.Second
	prt("Waiting", wait, "before exiting")
	time.Sleep(wait)
} //                                                                        main

// prt is like fmt.Println() but prefixes each line with a 'DEMO' tag.
func prt(args ...interface{}) {
	fmt.Println(append([]interface{}{"-------------> DEMO"}, args...)...)
} //                                                                         prt

// end
