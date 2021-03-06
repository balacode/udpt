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
}

// matchError retruns true if err contains the specified error description.
func matchError(err error, msg string) bool {
	if err == nil && (msg == "" || msg == "nil" || msg == "<nil>") {
		return true
	}
	return err != nil && strings.Contains(err.Error(), msg)
}

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
	failSetReadDeadline  bool
	failSetWriteBuffer   bool
	failSetWriteDeadline bool
	failReadFrom         bool
	failWrite            bool
	failWriteTo          bool
	failClose            bool
	//
	nSetReadDeadline  int
	nSetWriteBuffer   int
	nSetWriteDeadline int
	nReadFrom         int
	nWrite            int
	nWriteTo          int
	nClose            int
	//
	sertWriteBufferArg int
	writeDeadline      time.Time
	written            []byte
	readFromData       []byte
}

// ReadFrom implements PacketConn.ReadFrom().
func (mk *mockNetUDPConn) ReadFrom(b []byte) (int, net.Addr, error) {
	mk.nReadFrom++
	if mk.failReadFrom {
		return 0, nil, makeError(0xED19BF, "failed SetReadDeadline")
	}
	n := len(b)
	if len(mk.readFromData) > 0 {
		copy(b, mk.readFromData)
		n = len(mk.readFromData)
	}
	addr := &mockNetAddr{network: "udp", addr: "127.8.9.10:11"}
	return n, addr, nil
}

// Write implements Conn.Write().
func (mk *mockNetUDPConn) Write(b []byte) (int, error) {
	mk.nWrite++
	if mk.failWrite {
		return 0, makeError(0xEC15F0, "failed Write")
	}
	mk.written = append(mk.written, b...)
	return len(b), nil
}

// WriteTo implements PacketConn.WriteTo().
func (mk *mockNetUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	mk.nWriteTo++
	if mk.failWriteTo {
		return 0, makeError(0xEE40E7, "failed WriteTo")
	}
	mk.written = append(mk.written, b...)
	return len(b), nil
}

// SetReadDeadline implements Conn.SetReadDeadline().
func (mk *mockNetUDPConn) SetReadDeadline(time.Time) error {
	mk.nSetReadDeadline++
	if mk.failSetReadDeadline {
		return makeError(0xED5A2C, "failed SetReadDeadline")
	}
	return nil
}

// SetWriteBuffer sets the size of the transmit buffer of the connection.
func (mk *mockNetUDPConn) SetWriteBuffer(bytes int) error {
	mk.nSetWriteBuffer++
	if mk.failSetWriteBuffer {
		return makeError(0xE3EE33, "failed SetWriteBuffer")
	}
	mk.sertWriteBufferArg = bytes
	return nil
}

// SetWriteDeadline implements Conn.SetWriteDeadline().
func (mk *mockNetUDPConn) SetWriteDeadline(deadline time.Time) error {
	mk.nSetWriteDeadline++
	if mk.failSetWriteDeadline {
		return makeError(0xE63B56, "failed SetWriteDeadline")
	}
	mk.writeDeadline = deadline
	return nil
}

// Close closes the connection.
func (mk *mockNetUDPConn) Close() error {
	mk.nClose++
	if mk.failClose {
		return makeError(0xE60D82, "failed Close")
	}
	return nil
}

// -----------------------------------------------------------------------------

// mockWriteCloser is a mock io.WriteCloser with methods you can make fail.
type mockWriteCloser struct {
	failWrite bool
	failClose bool
}

// Write is a method of mockWriteCloser implementing io.WriteCloser.
//
// You can make it return an error by setting mockWriteCloser.failWrite.
//
func (mk *mockWriteCloser) Write(p []byte) (n int, err error) {
	if mk.failWrite {
		return 0, makeError(0xEA8F84, "failed mockWriteCloser.Write")
	}
	return len(p), nil
}

// Close is a method of mockWriteCloser implementing io.WriteCloser.
//
// You can make it return an error by setting mockWriteCloser.failClose.
//
func (mk *mockWriteCloser) Close() error {
	if mk.failClose {
		return makeError(0xEC5E59, "failed mockWriteCloser.Close")
	}
	return nil
}

// -----------------------------------------------------------------------------

// go test -run Test_Temp_
//
func Test_Temp_(t *testing.T) {
}

// end
