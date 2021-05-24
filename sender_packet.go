// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[sender_packet.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"io"
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
func (pk *senderPacket) IsDelivered() bool {
	ret := pk.confirmedHash != nil &&
		bytes.Equal(pk.sentHash, pk.confirmedHash)
	return ret
} //                                                                 IsDelivered

// Send encrypts and sends this packet through connection 'conn'.
func (pk *senderPacket) Send(conn netUDPConn, cipher SymmetricCipher) error {
	if conn == nil {
		return makeError(0xE4B1BA, "nil connection")
	}
	if cipher == nil {
		return makeError(0xE44F2A, "nil cipher")
	}
	ciphertext, err := cipher.Encrypt(pk.data)
	if err != nil {
		return makeError(0xEB39C3, err)
	}
	pk.sentTime = time.Now()
	_, err = io.Copy(conn, bytes.NewReader(ciphertext))
	if err != nil {
		return makeError(0xE93D1F, err)
	}
	return nil
} //                                                                        Send

// end
