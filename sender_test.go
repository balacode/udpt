// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[sender_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_Sender_*

// -----------------------------------------------------------------------------
// # Main Methods (ob *Sender)

// (ob *Sender) Send(name string, data []byte) error
//
// go test -run Test_Sender_Send_
//
func Test_Sender_Send_(t *testing.T) {
	{
		var snd Sender
		err := snd.Send("", nil)
		if snd.Config == nil {
			t.Error("0xE22B60")
		}
		if snd.Config.Cipher == nil {
			t.Error("0xEB62B4")
		}
		if !matchError(err, "invalid Sender.CryptoKey") {
			t.Error("0xE5BB36")
		}
		snd.CryptoKey = []byte("12345678901234567890123456789012")
		//
		snd.Config.Cipher = nil
		err = snd.Send("", nil)
		if !matchError(err, "nil Sender.Config.Cipher") {
			t.Error("0xE32EC6")
		}
		snd.Config.Cipher = &aesCipher{}
		//
		snd.Config.PacketSizeLimit = 65536
		err = snd.Send("", nil)
		if !matchError(err, "Sender.Config") {
			t.Error("0xE08E7C")
		}
		snd.Config.PacketSizeLimit = 65535 - 8
		//
		snd.Address = ""
		err = snd.Send("", nil)
		if !matchError(err, "Sender.Address") {
			t.Error("0xEC20C3")
		}
		snd.Address = "127.0.0.0"
		//
		snd.Port = 0
		err = snd.Send("", nil)
		if !matchError(err, "Sender.Port") {
			t.Error("0xE24E74")
		}
		snd.Port = 9876
		//
		snd.Config.VerboseSender = true
		snd.Config.SendRetries = 2
		snd.Config.ReplyTimeout = 1 * time.Second
		snd.Config.WriteTimeout = 1 * time.Second
		err = snd.Send("", nil)
		if !matchError(err, "undelivered packets") {
			t.Error("0xEB8B96")
		}
	}
} //                                                           Test_Sender_Send_

// -----------------------------------------------------------------------------
// # Informatory Properties (ob *Sender)

// (ob *Sender) AverageResponseMs() float64
//
// go test -run Test_Sender_AverageResponseMs_
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
// go test -run Test_Sender_TransferSpeedKBpS_
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
// go test -run Test_Sender_LogStats_
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
// go test -run Test_Sender_close_
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

// (ob *Sender) getPacketCount(length int) int
//
// go test -run Test_Sender_getPacketCount_
//
func Test_Sender_getPacketCount_(t *testing.T) {
	var ob Sender
	ob.Config = NewDefaultConfig()
	ob.Config.PacketPayloadSize = 1000
	//
	if ob.getPacketCount(0) != 0 {
		t.Error("0xE6C4D4")
	}
	if ob.getPacketCount(-100000) != 0 {
		t.Error("0xE81D08")
	}
	if ob.getPacketCount(1) != 1 {
		t.Error("0xE55EB9")
	}
	if ob.getPacketCount(1000) != 1 {
		t.Error("0xE87CB1")
	}
	if ob.getPacketCount(1001) != 2 {
		t.Error("0xE25DD0")
	}
	if ob.getPacketCount(10000) != 10 {
		t.Error("0xEE5EF4")
	}
} //                                                 Test_Sender_getPacketCount_

// (ob *Sender) logError(id uint32, args ...interface{}) error
//
// go test -run Test_Sender_logError_
//
func Test_Sender_logError_(t *testing.T) {
	var msg string
	var ob Sender
	ob.Config = NewDefaultConfig()
	ob.Config.LogFunc = func(args ...interface{}) {
		msg = fmt.Sprintln(args...)
	}
	ob.logError(0xE12345, "abc", 123)
	if msg != "ERROR 0xE12345: abc 123\n" {
		t.Error("0xE5CB5D")
	}
} //                                                       Test_Sender_logError_

// (ob *Sender) makePacket(data []byte) (*senderPacket, error)
//
// go test -run Test_Sender_makePacket_
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
		if !matchError(err, "PacketSizeLimit") {
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

// end
