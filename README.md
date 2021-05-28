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
        // provides existing data items for hashing by the Receiver. Only the
        // hash will be sent back to the sender, to confirm the transfer.
        ProvideData: func(k string) ([]byte, error) {
            fmt.Println("Receiver.ProvideData (k:" + k + "): " + received)
            return []byte(received), nil
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
```

## Security Notice:
This is a new project and its use of cryptography has not been reviewed by experts. While I make use of established crypto algorithms available in the standard Go library and would not "roll my own" encryption, there may be weaknesses in my application of the algorithms. Please use caution and do your own security asessment of the code. At present, this library uses AES-256 in Galois Counter Mode to encrypt each packet of data, including its headers, and SHA-256 for hashing binary resources that are being transferred.

## Version History:
This project is in its DRAFT stage: very unstable. At this point it works, but the API may change rapidly.

## Ideas:
- Create a drop-in replacement for TCP and TLS connections
- Implement some form of transfer control
- Improve performance
