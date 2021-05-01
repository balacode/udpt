// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                        /[connect.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
)

// connect connects to the Receiver specified by 'address'.
func connect(address string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, logError(0xEC7C6B, "(ResolveUDPAddr):", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr) // (*net.UDPConn, error)
	if err != nil {
		return nil, logError(0xE15CE1, "(DialUDP):", err)
	}
	// TODO: add this to ConfigSettings
	err = conn.SetWriteBuffer(16 * 1024 * 2014) // 16 MiB
	if err != nil {
		return nil, logError(0xE5F9C7, "(SetWriteBuffer):", err)
	}
	return conn, nil
} //                                                                     connect

// end
