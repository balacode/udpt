// -----------------------------------------------------------------------------
// github.com/balacode/udpt                         /[request_data_item_hash.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"encoding/hex"
)

// requestDataItemHash requests and waits for the listening receiver
// to return the hash of the data item named by 'name'. If the receiver
// can locate the data item, returns its hash, otherwise returns nil.
func requestDataItemHash(name string) []byte {
	err := Config.Validate()
	if err != nil {
		_ = logError(0xE5BC2E, err)
		return nil
	}
	conn, err := connect()
	if err != nil {
		_ = logError(0xE7DF8B, "(connect):", err)
		return nil
	}
	packet, err := NewPacket([]byte(DATA_ITEM_HASH + name))
	if err != nil {
		_ = logError(0xE1F8C5, "(NewPacket):", err)
		return nil
	}
	err = sendPacket(packet, conn) // *Packet, *net.UDPConn
	if err != nil {
		_ = logError(0xE7F316, "(sendPacket):", err)
		return nil
	}
	encryptedReply := make([]byte, Config.PacketSizeLimit)
	nRead, _ /*addr*/, err := readFromUDPConn(conn, encryptedReply)
	if err != nil {
		_ = logError(0xE97FC3, "(ReadFrom):", err)
		return nil
	}
	reply, err := aesDecrypt(encryptedReply[:nRead], Config.AESKey)
	if err != nil {
		_ = logError(0xE2B5A1, "(aesDecrypt):", err)
		return nil
	}
	var hash []byte
	if len(reply) > 0 {
		if !bytes.HasPrefix(reply, []byte(DATA_ITEM_HASH)) {
			_ = logError(0xE08AD4, ": invalid reply:", reply)
			return nil
		}
		hexHash := string(reply[len(DATA_ITEM_HASH):])
		if hexHash == "not_found" {
			return nil
		}
		hash, err = hex.DecodeString(hexHash)
		if err != nil {
			_ = logError(0xE5A4E7, "(hex.DecodeString):", err)
			return nil
		}
	}
	return hash
} //                                                         requestDataItemHash

// end
