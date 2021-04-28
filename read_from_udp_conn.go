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
	// contents of 'buf' is overwritten after every ReadFrom
	nRead, addr, err = conn.ReadFrom(tempBuf)
	if err != nil &&
		strings.Contains(err.Error(), "closed network connection") {
		err = nil
	}
	if err != nil {
		return 0, nil, logError(0xE0E0B1, "(ReadFrom):", err)
	}
	return nRead, addr, err
} //                                                             readFromUDPConn

// end
