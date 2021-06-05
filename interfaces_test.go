// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                /[interfaces_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"net"
	"testing"
)

// netDialUDP(network string, laddr, raddr *net.UDPAddr) (netUDPConn, error)
//
// go test -run Test_netDialUDP_*

// must succeed
func Test_netDialUDP_1(t *testing.T) {
	conn, err := netDialUDP("udp", nil,
		&net.UDPAddr{IP: []byte{127, 0, 0, 0}, Port: 9876})
	if conn == nil {
		t.Error("0xE3E28A")
	}
	if v, ok := conn.(*net.UDPConn); ok {
		if v == nil {
			t.Error("0xE9F3B4")
		}
	} else {
		t.Error("0xEE14A6")
	}
	if err != nil {
		t.Error("0xEF38D1")
	}
}

// must fail and return nil and an error because network is invalid
func Test_netDialUDP_2(t *testing.T) {
	conn, err := netDialUDP("badnet", nil, nil)
	if conn != nil {
		t.Error("0xE0E64E")
	}
	if !matchError(err, "unknown network") {
		t.Error("0xE75C27", "wrong error:", err)
	}
}

// must fail and return nil and an error because addres is not specified
func Test_netDialUDP_3(t *testing.T) {
	conn, err := netDialUDP("udp", nil, nil)
	if conn != nil {
		t.Error("0xE37D46")
	}
	if !matchError(err, "missing address") {
		t.Error("0xEB08F4", "wrong error:", err)
	}
}

// end
