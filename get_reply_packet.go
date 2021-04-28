// -----------------------------------------------------------------------------
// github.com/balacode/udpt                               /[get_reply_packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
)

// getReplyPacket _ _
func getReplyPacket(conn *net.UDPConn) []byte {
	err := Config.Validate()
	if err != nil {
		logError(0xE5BC2E, err)
		return []byte{}
	}
	encryptedReply := make([]byte, Config.PacketSizeLimit)
	nRead, _ /*addr*/, err := readFromUDPConn(conn, encryptedReply)
	if err != nil {
		logError(0xE97FC3, "(ReadFrom):", err)
		return []byte{}
	}
	ret, err := aesDecrypt(encryptedReply[:nRead], Config.AESKey)
	if err != nil {
		logError(0xE2B5A1, "(aesDecrypt):", err)
		return []byte{}
	}
	return ret
} //                                                              getReplyPacket

// end
