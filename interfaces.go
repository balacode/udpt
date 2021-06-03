// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                     /[interfaces.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"time"
)

// netUDPConn specifies the interface of net.UDPConn as used in this package.
type netUDPConn interface {

	// ReadFrom implements PacketConn.ReadFrom().
	ReadFrom([]byte) (int, net.Addr, error)

	// Write implements Conn.Write().
	Write(p []byte) (n int, err error)

	// WriteTo implements PacketConn.WriteTo().
	WriteTo(b []byte, addr net.Addr) (int, error)

	// SetReadDeadline implements Conn.SetReadDeadline().
	SetReadDeadline(time.Time) error

	// SetWriteBuffer sets the size of the transmit buffer of the connection.
	SetWriteBuffer(bytes int) error

	// SetWriteDeadline implements Conn.SetWriteDeadline().
	SetWriteDeadline(time.Time) error

	// Close closes the connection.
	Close() error
} //                                                                  netUDPConn

// netDialUDP wraps net.DialUDP and returns netUDPConn instead of *net.UDPConn
func netDialUDP(network string, laddr, raddr *net.UDPAddr) (netUDPConn, error) {
	return net.DialUDP(network, laddr, raddr)
} //                                                                  netDialUDP

// end
