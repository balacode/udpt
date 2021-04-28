// -----------------------------------------------------------------------------
// github.com/balacode/udpt                             /[read_from_udp_conn.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"strings"
	"time"
)

// readFromUDPConn reads data from the UDP connection 'conn'.
//
// 'tempBuf' contains a temporary buffer that holds the received
// packet's data. It is reused between calls to this function to
// avoid unneccessary memory allocations and de-allocations.
// The size of 'tempBuf' must be Config.PacketSizeLimit or greater.
//
func readFromUDPConn(
	conn *net.UDPConn,
	tempBuf []byte,
) (
	nRead int,
	addr net.Addr,
	err error,
) {
	err = conn.SetReadDeadline(time.Now().Add(Config.ReplyTimeout))
	if err != nil {
		return 0, nil, logError(0xE14A90, "(SetReadDeadline):", err)
	}
	// contents of 'tempBuf' is overwritten after every ReadFrom
	nRead, addr, err = conn.ReadFrom(tempBuf)
	if err != nil &&
		strings.Contains(err.Error(), "closed network connection") {
		err = nil
	}
	if err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			err = nil
		} else {
			err = logError(0xE0E0B1, "(ReadFrom):", err)
		}
	}
	return nRead, addr, err
} //                                                             readFromUDPConn

// end
