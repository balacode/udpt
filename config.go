// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[config.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"io"
	"os"
	"time"
)

// Config contains UDP and other default configuration settings.
// These settings normally don't need to be changed.
type Configuration struct {

	// -------------------------------------------------------------------------
	// Components:

	// Cipher is the object that handles encryption and decryption.
	//
	// It must implement the SymmetricCipher interface which is defined
	// in this package. If you don't specify Cipher, then encryption will
	// be done using the default AES-256 cipher used in this package.
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

	// LogWriter is the writer to which logError() and logInfo() output.
	// If you leave it nil, no logging will be done.
	LogWriter io.Writer

	// VerboseReceiver specifies if Receiver should
	// write informational log messages to LogWriter.
	VerboseReceiver bool

	// VerboseSender specifies if Sender should write
	// informational log messages to LogWriter.
	VerboseSender bool
} //                                                               Configuration

// NewDebugConfig returns configuration settings for debugging.
//
// You can specify an optional writer used for logging. If you omit it,
// logError() and logInfo() will print to standard output.
//
func NewDebugConfig(logWriter ...io.Writer) *Configuration {
	cf := NewDefaultConfig()
	cf.VerboseSender = true
	cf.VerboseReceiver = true
	if len(logWriter) > 0 {
		cf.LogWriter = logWriter[0]
	} else {
		cf.LogWriter = os.Stdout
	}
	return cf
} //                                                              NewDebugConfig

// NewDefaultConfig returns default configuration settings.
func NewDefaultConfig() *Configuration {
	return &Configuration{
		//
		// Components:
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
		ReplyTimeout:       10 * time.Second,
		SendPacketInterval: 1 * time.Millisecond,
		SendRetryInterval:  250 * time.Millisecond,
		SendWaitInterval:   25 * time.Millisecond,
		WriteTimeout:       10 * time.Second,
		//
		// Logging: (default nil/zero values)
	}
} //                                                            NewDefaultConfig

// Validate checks if all configuration parameters
// are set within acceptable limits.
//
// Returns nil if there is no problem, or the error instance.
//
func (cf *Configuration) Validate() error {
	//
	// Components:
	if cf.Cipher == nil {
		return makeError(0xE16FB9, "nil Configuration.Cipher")
	}
	if cf.Compressor == nil {
		return makeError(0xE5B3C1, "nil Configuration.Compressor")
	}
	// Limits:
	n := cf.PacketSizeLimit
	if n < 8 || n > (65535-8) {
		return makeError(0xE86C2A,
			"invalid Configuration.PacketSizeLimit:", n)
	}
	n = cf.PacketPayloadSize
	if n < 1 || n > (cf.PacketSizeLimit-200) {
		return makeError(0xE54BF4,
			"invalid Configuration.PacketPayloadSize:", n)
	}
	n = cf.SendBufferSize
	if n < 0 {
		return makeError(0xE27C2B,
			"invalid Configuration.SendBufferSize:", n)
	}
	n = cf.SendRetries
	if n < 0 {
		return makeError(0xE47C83,
			"invalid Configuration.SendRetries:", n)
	}
	return nil
} //                                                                    Validate

// end
