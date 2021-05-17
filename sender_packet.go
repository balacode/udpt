// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[sender_packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"io"
	"net"
	"time"
)

// senderPacket contains data, hash and timing details of
// a UDP packet (datagram) being sent by the Sender.
type senderPacket struct {
	data          []byte
	sentHash      []byte
	sentTime      time.Time
	confirmedHash []byte
	confirmedTime time.Time
} //                                                                senderPacket

// IsDelivered returns true if this packet has been successfully
// delivered (by receiving a successful confirmation packet).
func (ob *senderPacket) IsDelivered() bool {
	ret := bytes.Equal(ob.sentHash, ob.confirmedHash)
	return ret
} //                                                                 IsDelivered

// Send encrypts and sends this packet through connection 'conn'.
func (ob *senderPacket) Send(conn *net.UDPConn, cipher SymmetricCipher) error {
	if conn == nil {
		return makeError(0xE4B1BA, EInvalidArg, "nil conn")
	}
	if cipher == nil {
		return makeError(0xE44F2A, EInvalidArg, "nil cipher")
	}
	ciphertext, err := cipher.Encrypt(ob.data)
	if err != nil {
		return makeError(0xEB39C3, err)
	}
	ob.sentTime = time.Now()
	_, err = io.Copy(conn, bytes.NewReader(ciphertext))
	if err != nil {
		return makeError(0xE93D1F, err)
	}
	return nil
} //                                                                        Send

// end
