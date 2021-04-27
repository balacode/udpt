// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[config.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ConfigStruct contains UDP configuration settings, including
// the address of the server to which Send() should connect.
type ConfigStruct struct {

	// -------------------------------------------------------------------------
	// Main:

	// Address is the domain name or IP address of the listening receiver,
	// excluding the port number.
	Address string

	// Port is the port number of the listening server.
	// This number must be between 1 and 65535.
	Port int

	// AESKey the secret AES encryption key that must be shared
	// between the sendder (client) and the receiver (server).
	// This key must be exactly 32 bytes long.
	// This is the key AES uses for symmetric encryption.
	AESKey []byte

	// -------------------------------------------------------------------------
	// Limits:

	// PacketSizeLimit is the maximum size of a datagram in bytes,
	// including the headers, metadata and data payload.
	//
	// Maximum Transmission Unit (MTU):
	//
	// Internet Protocol requires hosts to process IP datagrams
	// of at least 576 bytes for IPv4 (or 1280 bytes for IPv6).
	// The IPv4 header is 20 bytes (or up to 60 with options).
	// UDP header is 8 bytes. 576 - 60 - 8 = 508 bytes available.
	//
	// The maximum Ethernet (v2) frame size is 1518 bytes, 18
	// of which are overhead, giving a usable size of 1500.
	// (To be on the safe side, we further reduce this by 50 bytes)
	//
	PacketSizeLimit int

	// PacketPayloadSize is the size of a single packet's payload, in bytes.
	// That is the part of the packet that contains actual useful data.
	// PacketPayloadSize must always be smaller that PacketSizeLimit.
	PacketPayloadSize int

	// -------------------------------------------------------------------------
	// Timeouts:

	// ReplyTimeout is the maximum time to wait for reply
	// datagram(s) to arrive in a UDP connection.
	ReplyTimeout time.Duration

	// SendRetries is the number of times for
	// Send() to retry sending lost packets.
	SendRetries int

	// WriteTimeout is the maximum time to
	// wait for writing to a UDP connection.
	WriteTimeout time.Duration

	// -------------------------------------------------------------------------
	// Logging:

	// VerboseReceiver specifies if the receiver should print
	// informational messages to the standard output.
	VerboseReceiver bool

	// VerboseSender specifies if Send() should print
	// informational messages to the standard output.
	VerboseSender bool
} //                                                                ConfigStruct

// Config contains global configuration settings.
var Config = ConfigStruct{

	// Main:
	Address: "127.0.0.1",
	Port:    0,
	AESKey:  []byte{},

	// Limits:
	PacketSizeLimit:   1450,
	PacketPayloadSize: 1024,

	// Timeouts:
	ReplyTimeout: 15 * time.Second,
	SendRetries:  10,
	WriteTimeout: 15 * time.Second,

	// Logging:
	VerboseReceiver: false,
	VerboseSender:   false,
} //                                                                      Config

// Validate checks the configuration to make sure all required fields like
// Address, Port and AESKey have been specified and are consistent.
//
// Returns nil if there is no problem, or the error code of the erorr.
//
func (ob *ConfigStruct) Validate() error {
	//
	// Main:
	if strings.TrimSpace(ob.Address) == "" {
		return errors.New("missing Config.Address")
	}
	if ob.Port < 1 || ob.Port > 65535 {
		return fmt.Errorf("invalid Config.Port: %d", ob.Port)
	}
	if len(ob.AESKey) != 32 {
		return fmt.Errorf(
			"Config.AESKey must be 32, but is %d bytes long",
			len(ob.AESKey),
		)
	}
	// Limits:
	n := ob.PacketSizeLimit
	if n < 8 || n > (65535-8) {
		return fmt.Errorf("invalid PacketSizeLimit: %d", n)
	}
	n = ob.PacketPayloadSize
	if n < 1 || n > (ob.PacketSizeLimit-200) {
		return fmt.Errorf("invalid PacketPayloadSize: %d", n)
	}
	return nil
} //                                                                    Validate

// end
