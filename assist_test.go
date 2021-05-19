// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[assist_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
)

// makeTestConn creates a UDP connection for testing.
func makeTestConn() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		panic("0xEE52A7")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic("0xE1E9E7")
	}
	return conn
} //                                                                makeTestConn
// end
