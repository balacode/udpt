// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                        /[connect.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"net"
	"strings"
)

/// add Address and Port args
// connect connects to the Receiver specified
// by Config.Address and Config.Port
func connect() (*net.UDPConn, error) {
	err := Config.Validate()
	if err != nil {
		return nil, logError(0xE5D78C, err)
	}
	addr := fmt.Sprintf("%s:%d",
		strings.TrimSpace(Config.Address), Config.Port,
	)
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, logError(0xEC7C6B, "(ResolveUDPAddr):", err)
	}
	conn, err := net.DialUDP("udp", nil, raddr) // (*net.UDPConn, error)
	if err != nil {
		return nil, logError(0xE15CE1, "(DialUDP):", err)
	}
	err = conn.SetWriteBuffer(16 * 1024 * 2014) // 16 MiB
	if err != nil {
		return nil, logError(0xE5F9C7, "(SetWriteBuffer):", err)
	}
	return conn, nil
} //                                                                     connect

// end
