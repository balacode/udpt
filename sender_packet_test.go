// -----------------------------------------------------------------------------
// github.com/balacode/udpt                             /[sender_packet_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"testing"
)

// to run all tests in this file:
// go test -v -run Test_senderPacket_*

// -----------------------------------------------------------------------------

// (pk *senderPacket) IsDelivered() bool
//
// go test -run Test_senderPacket_IsDelivered_
//
func Test_senderPacket_IsDelivered_(t *testing.T) {
	var pk senderPacket
	if pk.IsDelivered() != false {
		t.Error("0xEE17FE")
	}
	pk.sentHash = getHash([]byte("abc"))
	if pk.IsDelivered() != false {
		t.Error("0xE72B22")
	}
	pk.confirmedHash = getHash([]byte("abc"))
	if pk.IsDelivered() != true {
		t.Error("0xEE46BB")
	}
} //                                              Test_senderPacket_IsDelivered_

// (pk *senderPacket) Send(conn *net.UDPConn, cipher SymmetricCipher) error
//
// go test -run Test_senderPacket_Send_
//
func Test_senderPacket_Send_(t *testing.T) {
	{
		var pk senderPacket
		err := pk.Send(nil, nil)
		if !matchError(err, "nil conn") {
			t.Error("0xE31FF5")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		err := pk.Send(conn, nil)
		if !matchError(err, "nil cipher") {
			t.Error("0xE03CD3")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		cipher := &aesCipher{}
		err := pk.Send(conn, cipher)
		if !matchError(err, "key must be 32 bytes long") {
			t.Error("0xE12AB8")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		cipher := &aesCipher{cryptoKey: []byte{1, 2, 3}}
		err := pk.Send(conn, cipher)
		if !matchError(err, "key must be 32 bytes long") {
			t.Error("0xE53A3B")
		}
	}
	{
		var pk senderPacket
		conn := &net.UDPConn{}
		cipher := &aesCipher{}
		cipher.SetKey([]byte("12345678901234567890123456789012"))
		err := pk.Send(conn, cipher)
		if !matchError(err, "invalid argument") {
			t.Error("0xE65B73")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		cipher := &aesCipher{}
		cipher.SetKey([]byte("12345678901234567890123456789012"))
		err := pk.Send(conn, cipher)
		if err != nil {
			t.Error("0xED62D8")
		}
	}
} //                                                     Test_senderPacket_Send_

// end
