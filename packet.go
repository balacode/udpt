// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                         /[packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"io"
	"net"
	"time"
)

// Packet _ _
type Packet struct {
	data          []byte
	sentHash      []byte
	sentTime      time.Time
	confirmedHash []byte
	confirmedTime time.Time
} //                                                                      Packet

// isDelivered returns true if this packet has been successfully
// delivered (by receiving a successful confirmation packet).
func (ob *Packet) isDelivered() bool {
	ret := bytes.Equal(ob.sentHash, ob.confirmedHash)
	return ret
} //                                                                 isDelivered

// send encrypts and sends this packet through connection 'conn'.
func (ob *Packet) send(conn *net.UDPConn, cipher SymmetricCipher) error {
	if ob == nil {
		return makeError(0xE1D3B5, ENilReceiver)
	}
	if cipher == nil {
		return makeError(0xE54A9D, "nil cipher")
	}
	if conn == nil {
		return makeError(0xE4B1BA, EInvalidArg, ": conn is nil")
	}
	ciphertext, err := cipher.Encrypt(ob.data)
	if err != nil {
		return makeError(0xEB39C3, "(Encrypt):", err)
	}
	ob.sentTime = time.Now()
	_, err = io.Copy(conn, bytes.NewReader(ciphertext))
	if err != nil {
		return makeError(0xE93D1F, "(Copy):", err)
	}
	return nil
} //                                                                        send

// end
