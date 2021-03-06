// -----------------------------------------------------------------------------
// github.com/balacode/udpt                               /[read_and_decrypt.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"net"
	"strings"
	"time"
)

// errClosed error occurs when trying to read from a closed connection.
var errClosed = errors.New("use of closed network connection")

// errTimeout error occurs when a read operation times out.
//
// NOTE: this error is currently not checked for.
//
var errTimeout = errors.New("i/o timeout")

// -----------------------------------------------------------------------------

// readAndDecrypt reads data from the UDP connection 'conn'.
//
// 'tempBuf' contains a temporary buffer that holds the received
// packet's data. It is reused between calls to this function to
// avoid unnecessary memory allocations and de-allocations.
// The size of 'tempBuf' must be Config.PacketSizeLimit or greater.
//
func readAndDecrypt(
	conn netUDPConn,
	timeout time.Duration,
	decryptor SymmetricCipher,
	tempBuf []byte,
) (
	data []byte,
	addr net.Addr,
	err error,
) {
	if conn == nil {
		return nil, nil, makeError(0xE4ED27, "nil connection")
	}
	if decryptor == nil {
		return nil, nil, makeError(0xEF7F01, "nil decryptor")
	}
	if tempBuf == nil {
		return nil, nil, makeError(0xED80B0, "nil tempBuf")
	}
	dl := time.Now().Add(timeout)
	err = conn.SetReadDeadline(dl)
	if err != nil {
		return nil, nil, netError(err, 0xE09B6A)
	}
	// contents of 'tempBuf' is overwritten after every ReadFrom
	nRead, addr, err := conn.ReadFrom(tempBuf)
	if err != nil {
		return nil, nil, netError(err, 0xE0E0B1)
	}
	data, err = decryptor.Decrypt(tempBuf[:nRead])
	if err != nil {
		data, addr, err = nil, nil, makeError(0xE2B5A1, err)
	}
	return data, addr, err
} //                                                              readAndDecrypt

// netError filters out network errors for readAndDecrypt() and returns
// them as distinct error instances like errClosed and errTimeout.
//
// For other errors, it just calls makeError().
//
func netError(err error, otherErrorID uint32) error {
	if err == nil {
		return nil
	}
	errName := err.Error()
	switch {
	// don't log a closed connection or i/o timeout:
	// these are expected, so just return errClosed or errTimeout
	case strings.Contains(errName, errClosed.Error()):
		err = errClosed
	case strings.Contains(errName, errTimeout.Error()):
		err = errTimeout
	default:
		// log any other unexpected error here
		err = makeError(otherErrorID, err)
	}
	return err
} //                                                                    netError

// end
