## udpt
UDP Transport

[![Go Report Card](https://goreportcard.com/badge/github.com/balacode/udpt)](https://goreportcard.com/report/github.com/balacode/udpt)
[![godoc](https://godoc.org/github.com/balacode/udpt?status.svg)](https://godoc.org/github.com/balacode/udpt)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Compresses, encrypts and transfers data between a sender and receiver using UDP protocol.

## Features and Design Aims:
- Avoid the overhead of establishing a TCP or TCP+TLS handshake.
- Reliable transfer of data using an unreliable UDP connection.
- Uses AES-256 symmetric cipher for encryption.
- Uses zlib library for data compression.
- No third-party dependencies. Only uses the standard library.
- Readable, understandable code with explanatory comments.

## Installation:

```bash
    go get github.com/balacode/udpt
```

## Hello World:

This demo runs a receiver using RunReceiver() which listens for incoming data,
then sends a "Hello World" to the receiver using Sender.SendString().

```go
package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/balacode/udpt"
)

func main() {
    // the AES key shared by the sender and receiver: must be 32 bytes log
    var aesKey = []byte("aA2Xh41FiC4Wtj3e5b2LbytMdn6on7P0")
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
        Port:   1234,
        AESKey: aesKey,
        Config: cfg,
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
        Address: "127.0.0.1",
        Port:    1234,
        AESKey:  aesKey,
        Config:  cfg,
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
```

## Version History:
This project is in its DRAFT stage: very unstable. At this point it works, but the API may change rapidly.

## Ideas:
- Write unit tests
- Create a drop-in replacement for TCP and TLS connections
- Implement some form of transfer control
- Improve performance
- Allow multiple Senders and Receivers that use different Address and Port values.
