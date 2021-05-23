// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[assist_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"strings"
	"testing"
	"time"
)

// makeTestConn creates a UDP connection for testing.
func makeTestConn() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		panic(makeError(0xEE52A7, err).Error())
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(makeError(0xE1E9E7, err).Error())
	}
	return conn
} //                                                                makeTestConn

// matchError retruns true if err contains the specified error message.
func matchError(err error, msg string) bool {
	if err == nil && (msg == "" || msg == "nil" || msg == "<nil>") {
		return true
	}
	return err != nil && strings.Contains(err.Error(), msg)
} //                                                                  matchError

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

// -----------------------------------------------------------------------------

// mockWriteCloser is a mock io.WriteCloser with methods you can make fail.
type mockWriteCloser struct {
	failWrite bool
	failClose bool
} //                                                             mockWriteCloser

// Write is a method of mockWriteCloser implementing io.WriteCloser.
//
// You can make it return an error by setting mockWriteCloser.failWrite.
//
func (mk *mockWriteCloser) Write(p []byte) (n int, err error) {
	if mk.failWrite {
		return 0, makeError(0xEA8F84, "failed mockWriteCloser.Write")
	}
	return len(p), nil
} //                                                                       Write

// Close is a method of mockWriteCloser implementing io.WriteCloser.
//
// You can make it return an error by setting mockWriteCloser.failClose.
//
func (mk *mockWriteCloser) Close() error {
	if mk.failClose {
		return makeError(0xEC5E59, "failed mockWriteCloser.Close")
	}
	return nil
} //                                                                       Close

// -----------------------------------------------------------------------------

// go test -run Test_Temp_
//
func Test_Temp_(t *testing.T) {
} //                                                                  Test_Temp_

// end
