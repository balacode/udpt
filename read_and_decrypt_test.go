// -----------------------------------------------------------------------------
// github.com/balacode/udpt                          /[read_and_decrypt_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"net"
	"testing"
	"time"
)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// readAndDecrypt(conn netUDPConn, tempBuf []byte, timeout time.Duration)
//
// go test -run Test_readAndDecrypt_*

// must succeed
func Test_readAndDecrypt_1(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{},      // conn
		time.Second,            // timeout
		newTestAESCipher(t),    // decryptor
		newTestAESCiphertext(), // tempBuf
	)
	if string(data) != "abc" {
		t.Error("0xE6E2DD")
	}
	if addr == nil ||
		addr.Network() != "udp" || addr.String() != "127.8.9.10:11" {
		t.Error("0xE58EB3")
	}
	if err != nil {
		t.Error("0xEA21C0", err)
	}
} //                                                       Test_readAndDecrypt_1

// must fail because connection is nil, before any other checks
func Test_readAndDecrypt_2(t *testing.T) {
	data, addr, err := readAndDecrypt(
		nil,         // conn <- failure
		time.Second, // timeout
		nil,         // decryptor
		nil,         // tempBuf
	)
	if len(data) != 0 {
		t.Error("0xEF5EA3")
	}
	if addr != nil {
		t.Error("0xE6FD31")
	}
	if !matchError(err, "nil connection") {
		t.Error("0xE7B3EA", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_2

// must fail because decryptor is nil
func Test_readAndDecrypt_3(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{},   // conn,
		time.Second,         // timeout
		nil,                 // decryptor <-failure: nil decryptor
		make([]byte, 65536), // tempBuf
	)
	if len(data) != 0 {
		t.Error("0xEA7CA1")
	}
	if addr != nil {
		t.Error("0xED4D55")
	}
	if !matchError(err, "nil decryptor") {
		t.Error("0xED13F3", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_3

// must fail because tempBuf is nil
func Test_readAndDecrypt_4(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{},   // conn,
		time.Second,         // timeout
		newTestAESCipher(t), // decryptor
		nil,                 // tempBuf <-failure: nil tempBuf
	)
	if len(data) != 0 {
		t.Error("0xE6E42A")
	}
	if addr != nil {
		t.Error("0xE65A6E")
	}
	if !matchError(err, "nil tempBuf") {
		t.Error("0xE12E41", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_4

// must fail when conn.SetReadDeadline() fails
func Test_readAndDecrypt_5(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{failSetReadDeadline: true}, // conn <-failure
		time.Second,         // timeout
		newTestAESCipher(t), // decryptor
		make([]byte, 65536), // tempBuf
	)
	if len(data) != 0 {
		t.Error("0xEC3C98")
	}
	if addr != nil {
		t.Error("0xED8ED3")
	}
	if !matchError(err, "failed SetReadDeadline") {
		t.Error("0xEE04B0", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_5

// must fail when conn.ReadFrom() fails
func Test_readAndDecrypt_6(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{failReadFrom: true}, // conn <-failure
		time.Second,                         // timeout
		newTestAESCipher(t),                 // decryptor
		make([]byte, 65536),                 // tempBuf
	)
	if len(data) != 0 {
		t.Error("0xED40F6")
	}
	if addr != nil {
		t.Error("0xE14D91")
	}
	if !matchError(err, "failed SetReadDeadline") {
		t.Error("0xE3E57D", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_6

// must fail because ciphertext is garbage
func Test_readAndDecrypt_7(t *testing.T) {
	data, addr, err := readAndDecrypt(
		&mockNetUDPConn{},              // conn
		time.Second,                    // timeout
		newTestAESCipher(t),            // decryptor
		[]byte{0xA8, 0xE1, 0x7D, 0xD6}, // tempBuf <-failure: bad ciphertext
	)
	if data != nil {
		t.Error("0xEC9C89")
	}
	if addr != nil {
		t.Error("0xED11F4")
	}
	if !matchError(err, "invalid ciphertext") {
		t.Error("0xEA53B8", "wrong error:", err)
	}
} //                                                       Test_readAndDecrypt_7

// -----------------------------------------------------------------------------

// netError(err error, otherErrorID uint32) error
//
// go test -run Test_netError_
//
func Test_netError_(t *testing.T) {
	err := netError(nil, 0xE12345)
	if err != nil {
		t.Error("0xEA7F2E", err)
	}
	err = netError(errors.New("..use of closed network connection.."), 0xE47EB8)
	if err != errClosed {
		t.Error("0xE8BD57")
	}
	err = netError(errors.New("..i/o timeout.."), 0xEA11C7)
	if err != errTimeout {
		t.Error("0xE53ED7")
	}
	err = netError(errors.New("some other error"), 0xE3E3D8)
	if !matchError(err, "some other error") {
		t.Error("0xE6F1BB", "wrong error:", err)
	}
} //                                                              Test_netError_

// -----------------------------------------------------------------------------

// mockNetAddr is a mock net.Addr implementation which can
// be made to return the network and address you want.
type mockNetAddr struct {
	network string
	addr    string
}

// Network is the name of the network (e.g. "udp")
func (mk *mockNetAddr) Network() string { return mk.network }

// String form of the address (e.g. "127.0.0.1:9876")
func (mk *mockNetAddr) String() string { return mk.addr }

// -----------------------------------------------------------------------------

// mockNetUDPConn is a mock net.UDPConn with methods you can make fail.
type mockNetUDPConn struct {
	failSetReadDeadline bool
	failReadFrom        bool
	failClose           bool
} //                                                              mockNetUDPConn

func (mk *mockNetUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	if mk.failReadFrom {
		return 0, nil, makeError(0xED19BF, "failed SetReadDeadline")
	}
	addr := &mockNetAddr{network: "udp", addr: "127.8.9.10:11"}
	return len(b), addr, nil
} //                                                                    ReadFrom

func (mk *mockNetUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return 0, nil
} //                                                                     WriteTo

func (mk *mockNetUDPConn) SetReadDeadline(time.Time) error {
	if mk.failSetReadDeadline {
		return makeError(0xED5A2C, "failed SetReadDeadline")
	}
	return nil
} //                                                             SetReadDeadline

func (mk *mockNetUDPConn) SetWriteDeadline(time.Time) error {
	return nil
} //                                                            SetWriteDeadline

func (mk *mockNetUDPConn) Close() error {
	if mk.failClose {
		return makeError(0xE60D82, "failed Close")
	}
	return nil
} //                                                                       Close

// end
