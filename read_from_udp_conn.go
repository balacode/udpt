// -----------------------------------------------------------------------------
// github.com/balacode/udpt                             /[read_from_udp_conn.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"net"
	"strings"
	"time"
)

var (
	// ErrClosed error occurs when trying to read from a closed connection.
	ErrClosed = errors.New("use of closed network connection")

	// ErrTimeout error occurs when a read operation times out.
	ErrTimeout = errors.New("i/o timeout")
)

// readFromUDPConn reads data from the UDP connection 'conn'.
//
// 'tempBuf' contains a temporary buffer that holds the received
// packet's data. It is reused between calls to this function to
// avoid unnecessary memory allocations and de-allocations.
// The size of 'tempBuf' must be Config.PacketSizeLimit or greater.
//
func readFromUDPConn(
	conn *net.UDPConn,
	tempBuf []byte,
	timeout time.Duration,
) (
	nRead int,
	addr net.Addr,
	err error,
) {
	if conn == nil {
		return 0, nil, logError(0xE4ED27, EInvalidArg)
	}
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return 0, nil, logError(0xE14A90, "(SetReadDeadline):", err)
	}
	// contents of 'tempBuf' is overwritten after every ReadFrom
	nRead, addr, err = conn.ReadFrom(tempBuf)
	if err != nil {
		errName := err.Error()
		switch {
		// don't log a closed connection or i/o timeout:
		// these are expected, so just return ErrClosed or ErrTimeout
		case strings.Contains(errName, ErrClosed.Error()):
			err = ErrClosed
		case strings.Contains(errName, ErrTimeout.Error()):
			err = ErrTimeout
		default:
			// log any other unexpected error here
			err = logError(0xE0E0B1, "(readFromUDPConn):", err)
		}
	}
	return nRead, addr, err
} //                                                             readFromUDPConn

// end
