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
	//
	// Specify required configuration parameters:
	//
	// Address is used by the sender to connect to the Receiver.
	//
	// Port is the port number on which the sender
	// sends and the receiver listens.
	//
	// AESKey is the secret AES encryption key shared by the
	// Sender and the Receiver. It must be exactly 32 bytes.
	//
	udpt.Config.Address = "127.0.0.1"
	udpt.Config.Port = 1234
	udpt.Config.AESKey = []byte{
		0xC4, 0x53, 0x67, 0xA7, 0xB7, 0x94, 0xE5, 0x30,
		0x6C, 0x4F, 0x43, 0x6C, 0xA9, 0x33, 0x85, 0xEA,
		0x1C, 0x37, 0xE3, 0x66, 0x7F, 0x14, 0x05, 0xE6,
		0x2F, 0x8F, 0xC6, 0x12, 0x67, 0x04, 0x86, 0xD1,
	}
	// disable log caching and enable verbose messages.
	// This should be done only during demos/prototyping/debugging.
	udpt.LogBufferSize = -1
	udpt.Config.VerboseSender = true
	udpt.Config.VerboseReceiver = true
	//
	prt("Running the receiver")
	//
	// receiveData is the function that receives data sent to the receiver
	receiveData := func(name string, data []byte) error {
		div := strings.Repeat("##", 40)
		prt(div)
		prt("You should see a 'Hello World!' message below:")
		prt(div)
		prt("Receiver's receiveData()")
		prt("Received name:", name)
		prt("Received data:", string(data))
		prt(div)
		return nil
	}
	// provideData is the function used to read back the data previously
	// received by the receiver. This data is never sent back to the
	// sender. It is only used to generate a hash that is sent to
	// the sender only to confirm that a data item has been sent.
	provideData := func(name string) ([]byte, error) {
		prt("Receiver's provideData()")
		return nil, nil
	}
	udpt.RunReceiver(receiveData, provideData)
	//
	time.Sleep(1 * time.Second)
	prt("Sending a message")
	sender := udpt.Sender{Config: udpt.Config}
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
