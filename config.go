// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[config.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"time"
)

// Config contains UDP and other default configuration settings.
// These settings normally don't need to be changed.
type Configuration struct {

	// Cipher is the object that handles encryption and decryption.
	//
	// It must implement the SymmetricCipher interface which is defined in
	// this package. If you don't specify Cipher, then encryption will be done
	// using the default AES-256 cipher (aesCipher) used in this package.
	//
	Cipher SymmetricCipher

	// Compressor handles compression and uncompression.
	Compressor Compression

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

	// SendBufferSize is size of the write buffer used by Send(), in bytes.
	SendBufferSize int

	// SendRetries is the number of times for
	// Send() to retry sending lost packets.
	SendRetries int

	// -------------------------------------------------------------------------
	// Timeouts and Intervals:

	// ReplyTimeout is the maximum time to wait for reply
	// datagram(s) to arrive in a UDP connection.
	ReplyTimeout time.Duration

	// SendPacketInterval is the time to wait between sending packets.
	SendPacketInterval time.Duration

	// SendRetryInterval is the time for Sender.Send() to
	// wait before retrying to send undelivered packets.
	SendRetryInterval time.Duration

	// SendWaitInterval is the amount of time Sender() should sleep
	// in the loop, before checking if a confirmation has arrived.
	SendWaitInterval time.Duration

	// WriteTimeout is the maximum time to
	// wait for writing to a UDP connection.
	WriteTimeout time.Duration

	// -------------------------------------------------------------------------
	// Logging:

	// LogFunc is the function used to log logError() and logInfo() messages.
	// If you leave it nil, no logging will be done.
	LogFunc func(args ...interface{})

	// VerboseReceiver specifies if the receiver should print
	// informational messages to the standard output.
	VerboseReceiver bool

	// VerboseSender specifies if Send() should print
	// informational messages to the standard output.
	VerboseSender bool
} //                                                               Configuration

// NewDebugConfig returns configuration settings for debugging.
//
// You can specify an optional log function for logging.
// If you omit it, logError() and logInfo() output will
// use LogPrint, which just prints to standard output.
//
// Tip: to log output to specific file in addition to standard output, use:
//
// udpt.NewDebugConfig(udpt.MakeLogFunc(true, "your_file_name"))
//
// If you pass multiple arguments, only the first will be used.
//
func NewDebugConfig(logFunc ...func(args ...interface{})) *Configuration {
	cfg := NewDefaultConfig()
	cfg.VerboseSender = true
	cfg.VerboseReceiver = true
	//
	if len(logFunc) > 0 {
		const printMsg = true
		cfg.LogFunc = logFunc[0]
	} else {
		cfg.LogFunc = LogPrint
	}
	return cfg
} //                                                              NewDebugConfig

// NewDefaultConfig returns default configuration settings.
func NewDefaultConfig() *Configuration {
	return &Configuration{
		Cipher:     &aesCipher{},
		Compressor: &zlibCompressor{},
		//
		// Limits:
		PacketSizeLimit:   1450,
		PacketPayloadSize: 1024,
		SendBufferSize:    16 * 1024 * 2014, // 16 MiB
		SendRetries:       10,
		//
		// Timeouts and Intervals:
		ReplyTimeout:       15 * time.Second,
		SendPacketInterval: 2 * time.Millisecond,
		SendRetryInterval:  250 * time.Millisecond,
		SendWaitInterval:   50 * time.Millisecond,
		WriteTimeout:       15 * time.Second,
		//
		// Logging:
		LogFunc:         nil,
		VerboseReceiver: false,
		VerboseSender:   false,
	}
} //                                                            NewDefaultConfig

// Validate checks if all configuration parameters
// are set within acceptable limits.
//
// Returns nil if there is no problem, or the error value.
//
func (ob *Configuration) Validate() error {
	if ob == nil {
		return makeError(0xE8E9E5, ENilReceiver+" in Configuration")
	}
	if ob.Cipher == nil {
		return makeError(0xE5D4AB, "nil Configuration.Cipher")
	}
	if ob.Compressor == nil {
		return makeError(0xE5B3C1, "nil Configuration.Compressor")
	}
	n := ob.PacketSizeLimit
	if n < 8 || n > (65535-8) {
		return makeError(0xE86C2A,
			"invalid Configuration.PacketSizeLimit:", n)
	}
	n = ob.PacketPayloadSize
	if n < 1 || n > (ob.PacketSizeLimit-200) {
		return makeError(0xE54BF4,
			"invalid Configuration.PacketPayloadSize:", n)
	}
	n = ob.SendBufferSize
	if n < 0 {
		return makeError(0xE27C2B,
			"invalid Configuration.SendBufferSize:", n)
	}
	n = ob.SendRetries
	if n < 0 {
		return makeError(0xE1C8A6,
			"invalid Configuration.SendRetries:", n)
	}
	return nil
} //                                                                    Validate

// end
