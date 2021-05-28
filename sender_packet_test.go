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
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (pk *senderPacket) Send(conn *net.UDPConn, cipher SymmetricCipher) error
//
// go test -run Test_senderPacket_Send_*

// must succeed
func Test_senderPacket_Send_1(t *testing.T) {
	var pk senderPacket
	conn := makeTestConn()
	cipher := &aesCipher{}
	cipher.SetKey([]byte("12345678901234567890123456789012"))
	err := pk.Send(conn, cipher)
	if err != nil {
		t.Error("0xED62D8", err)
	}
}

// must fail when passed a nil connection
func Test_senderPacket_Send_2(t *testing.T) {
	var pk senderPacket
	err := pk.Send(nil, nil)
	if !matchError(err, "nil conn") {
		t.Error("0xE31FF5", "wrong error:", err)
	}
}

// must fail when passed an invalid connection
func Test_senderPacket_Send_3(t *testing.T) {
	var pk senderPacket
	conn := &net.UDPConn{} // bad connection
	cipher := &aesCipher{}
	cipher.SetKey([]byte("12345678901234567890123456789012"))
	err := pk.Send(conn, cipher)
	if !matchError(err, "invalid argument") {
		// TODO: above error description may differ on Linux or Mac OS
		t.Error("0xE65B73", "wrong error:", err)
	}
}

// must fail when passed a nil cipher
func Test_senderPacket_Send_4(t *testing.T) {
	var pk senderPacket
	conn := makeTestConn()
	err := pk.Send(conn, nil)
	if !matchError(err, "nil cipher") {
		t.Error("0xE03CD3", "wrong error:", err)
	}
}

// must fail when the encryption key is not set
func Test_senderPacket_Send_5(t *testing.T) {
	var pk senderPacket
	conn := makeTestConn()
	cipher := &aesCipher{}
	err := pk.Send(conn, cipher)
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE12AB8", "wrong error:", err)
	}
}

// must fail when the encryption key is invalid
func Test_senderPacket_Send_6(t *testing.T) {
	var pk senderPacket
	conn := makeTestConn()
	cipher := &aesCipher{cryptoKey: []byte{1, 2, 3}}
	err := pk.Send(conn, cipher)
	if !matchError(err, "AES-256 key must be 32 bytes long") {
		t.Error("0xE53A3B", "wrong error:", err)
	}
}

// end
