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
	ReadFrom([]byte) (int, net.Addr, error)
	Write(p []byte) (n int, err error)
	WriteTo(b []byte, addr net.Addr) (int, error)
	SetReadDeadline(time.Time) error
	SetWriteBuffer(bytes int) error
	SetWriteDeadline(time.Time) error
	Close() error
} //                                                                  netUDPConn

// end
