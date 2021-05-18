// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[sender_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run _Sender_

// -----------------------------------------------------------------------------
// # Informatory Properties (ob *Sender)

// (ob *Sender) AverageResponseMs() float64
//
// go test -run _Sender_AverageResponseMs_
//
func Test_Sender_AverageResponseMs_(t *testing.T) {
	var snd Sender
	if n := snd.AverageResponseMs(); n < 0 || n > 0 {
		t.Error("0xE29B40")
	}
	snd.stats.packetsDelivered = 10
	snd.stats.transferTime = time.Millisecond
	if n := snd.AverageResponseMs(); n < 0.1 || n > 0.1 {
		t.Error("0xE01B5F")
	}
} //                                              Test_Sender_AverageResponseMs_

// (ob *Sender) TransferSpeedKBpS() float64
//
// go test -run _Sender_TransferSpeedKBpS_
//
func Test_Sender_TransferSpeedKBpS_(t *testing.T) {
	var snd Sender
	if n := snd.TransferSpeedKBpS(); n < 0 || n > 0 {
		t.Error("0xEE99D3")
	}
	snd.stats.transferTime = time.Second
	snd.stats.bytesDelivered = 88 * 1024
	if n := snd.TransferSpeedKBpS(); n < 88 || n > 88 {
		t.Error("0xED0BD8")
	}
} //                                              Test_Sender_TransferSpeedKBpS_

// -----------------------------------------------------------------------------
// # Informatory Methods (ob *Sender)

// (ob *Sender) LogStats(logFunc ...interface{})
//
// go test -run _Sender_LogStats_
//
func Test_Sender_LogStats_(t *testing.T) {
	var sb strings.Builder
	fmtPrintln := func(v ...interface{}) (int, error) {
		sb.WriteString(fmt.Sprintln(v...))
		return 0, nil
	}
	logPrintln := func(v ...interface{}) {
		sb.WriteString(fmt.Sprintln(v...))
	}
	test := func(logFunc interface{}) {
		sb.Reset()
		//
		var snd Sender
		snd.Config = NewDefaultConfig()
		snd.Config.LogFunc = logPrintln
		snd.packets = []senderPacket{{
			sentHash:      []byte{0},
			confirmedHash: []byte{0},
		}}
		snd.LogStats(logFunc)
		//
		got := sb.String()
		expect := "" +
			"SN: 0    T0: 0001-01-01 00:00:00 +000 T1: NONE ✔ 0.0 ms\n" +
			"B. delivered: 0\n" +
			"Bytes lost  : 0\n" +
			"P. delivered: 0\n" +
			"Packets lost: 0\n" +
			"Time in item: 0.0 s\n" +
			"Avg./ Packet: 0.0 ms\n" +
			"Trans. speed: 0.0 KiB/s\n"
		//
		if got != expect {
			t.Error("\n" + "expect:\n" + expect + "\n" + "got:\n" + got)
			fmt.Println([]byte(expect))
			fmt.Println([]byte(got))
		}
	}
	test(nil)
	test(logPrintln)
	test(fmtPrintln)
} //                                                       Test_Sender_LogStats_

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (ob *Sender)

// (ob *Sender) close() error
//
// go test -run _Sender_closes_
//
func Test_Sender_close_(t *testing.T) {
	var snd Sender
	err := snd.close()
	if err != nil {
		t.Error("0xEE7E05")
	}
	snd.conn = makeTestConn()
	err = snd.conn.Close()
	if err != nil {
		t.Error("0xE3DD56")
	}
	err = snd.conn.Close()
	if err == nil {
		t.Error("0xE5FE16")
	}
} //                                                          Test_Sender_close_

// -----------------------------------------------------------------------------
// # Internal Helper Methods (ob *Sender)

// (ob *Sender) makePacket(data []byte) (*senderPacket, error)
//
// go test -run _Sender_makePacket_
//
func Test_Sender_makePacket_(t *testing.T) {
	{
		var snd Sender
		snd.Config = NewDefaultConfig()
		data := make([]byte, snd.Config.PacketSizeLimit+1)
		packet, err := snd.makePacket(data)
		if packet != nil {
			t.Error("0xE0FE30")
		}
		if err == nil {
			t.Error("0xEE76D9")
		} else if !strings.Contains(err.Error(), "PacketSizeLimit") {
			t.Error("0xE69EA5")
		}
	}
	{
		var snd Sender
		snd.Config = NewDefaultConfig()
		packet, err := snd.makePacket([]byte{1, 2, 3})
		if err != nil {
			t.Error("0xE0AE90")
		}
		expectData := []byte{1, 2, 3}
		if !bytes.Equal(packet.data, expectData) {
			t.Error("0xEF4E82")
		}
		expectSentHash := getHash([]byte{1, 2, 3})
		if !bytes.Equal(packet.sentHash, expectSentHash) {
			t.Error("0xE51B95")
		}
		n := time.Since(packet.sentTime)
		if n > time.Millisecond {
			t.Error("0xE1FA4B")
		}
		if packet.confirmedHash != nil {
			t.Error("0xEA0E4B")
		}
		if !packet.confirmedTime.IsZero() {
			t.Error("0xE21EB4")
		}
	}
} //                                                     Test_Sender_makePacket_

// -----------------------------------------------------------------------------

// makeTestConn creates a UDP connection for testing.
func makeTestConn() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9876")
	if err != nil {
		panic("0xEE52A7")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic("0xE1E9E7")
	}
	return conn
} //                                                                makeTestConn

// end
