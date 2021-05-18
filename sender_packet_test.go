// -----------------------------------------------------------------------------
// github.com/balacode/udpt                             /[sender_packet_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"strings"
	"testing"
)

// to run all tests in this file:
// go test -v -run _senderPacket_

// -----------------------------------------------------------------------------

// (ob *senderPacket) IsDelivered() bool
//
// go test -run _senderPacket_IsDelivered_
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

// (ob *senderPacket) Send(conn *net.UDPConn, cipher SymmetricCipher) error
//
// go test -run _senderPacket_Send_
//
func Test_senderPacket_Send_(t *testing.T) {
	{
		var pk senderPacket
		err := pk.Send(nil, nil)
		if err == nil {
			t.Error("0xE66F44")
		} else if !strings.Contains(err.Error(), "nil conn") {
			t.Error("0xE31FF5")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		err := pk.Send(conn, nil)
		if err == nil {
			t.Error("0xEE28A6")
		} else if !strings.Contains(err.Error(), "nil cipher") {
			t.Error("0xE03CD3")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		cipher := &aesCipher{}
		err := pk.Send(conn, cipher)
		if err == nil {
			t.Error("0xE55D1A")
		} else if !strings.Contains(err.Error(), "key must be 32 bytes long") {
			t.Error("0xE12AB8")
		}
	}
	{
		var pk senderPacket
		conn := makeTestConn()
		cipher := &aesCipher{cryptoKey: []byte{1, 2, 3}}
		err := pk.Send(conn, cipher)
		if err == nil {
			t.Error("0xE33D78")
		} else if !strings.Contains(err.Error(), "key must be 32 bytes long") {
			t.Error("0xE53A3B")
		}
	}
	{
		var pk senderPacket
		conn := &net.UDPConn{}
		cipher := &aesCipher{}
		cipher.SetKey([]byte("12345678901234567890123456789012"))
		err := pk.Send(conn, cipher)
		if err == nil {
			t.Error("0xEF5B42")
		} else if !strings.Contains(err.Error(), "invalid argument") {
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
