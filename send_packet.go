// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[send_packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"io"
	"net"
	"time"
)

// sendPacket encrypts and sends packet through connection conn.
// Returns an error if the packet could not be encrypted or sent.
func sendPacket(
	packet *Packet,
	aesKey []byte,
	conn *net.UDPConn,
) error {
	if packet == nil {
		return logError(0xE1D3B5, ENilReceiver)
	}
	if conn == nil {
		return logError(0xE4B1BA, ENilReceiver)
	}
	encryptedReq, err := aesEncrypt(packet.data, aesKey)
	if err != nil {
		return logError(0xEB39C3, "(aesEncrypt):", err)
	}
	packet.sentTime = time.Now()
	_, err = io.Copy(conn, bytes.NewReader(encryptedReq))
	if err != nil {
		return logError(0xE93D1F, "(Copy):", err)
	}
	return nil
} //                                                                  sendPacket

// end
