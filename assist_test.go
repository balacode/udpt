// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[assist_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"strings"
	"testing"
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
