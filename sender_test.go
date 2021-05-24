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
// go test -v -run Test_Sender_*

// -----------------------------------------------------------------------------
// # Main Methods (sd *Sender)

// (sd *Sender) Send(name string, data []byte) error
//
// go test -run Test_Sender_Send_
//
func Test_Sender_Send_(t *testing.T) {
	var sd Sender
	err := sd.Send("", nil)
	if sd.Config == nil {
		t.Error("0xE22B60")
	}
	if sd.Config.Cipher == nil {
		t.Error("0xEB62B4")
	}
	if !matchError(err, "invalid Sender.CryptoKey") {
		t.Error("0xE5BB36", "wrong error:", err)
	}
	sd.CryptoKey = []byte("12345678901234567890123456789012")
	//
	sd.Config.Cipher = nil
	err = sd.Send("", nil)
	if !matchError(err, "nil Sender.Config.Cipher") {
		t.Error("0xE32EC6", "wrong error:", err)
	}
	sd.Config.Cipher = &aesCipher{}
	//
	sd.Config.PacketSizeLimit = 65536
	err = sd.Send("", nil)
	if !matchError(err, "Sender.Config") {
		t.Error("0xE08E7C", "wrong error:", err)
	}
	sd.Config.PacketSizeLimit = 65535 - 8
	//
	sd.Address = ""
	err = sd.Send("", nil)
	if !matchError(err, "Sender.Address") {
		t.Error("0xEC20C3", "wrong error:", err)
	}
	sd.Address = "127.0.0.0"
	//
	sd.Port = 0
	err = sd.Send("", nil)
	if !matchError(err, "Sender.Port") {
		t.Error("0xE24E74", "wrong error:", err)
	}
	sd.Port = 9876
	//
	sd.Config.VerboseSender = true
	sd.Config.SendRetries = 2
	sd.Config.ReplyTimeout = 500 * time.Millisecond
	sd.Config.WriteTimeout = 500 * time.Millisecond
	err = sd.Send("", nil)
	if !matchError(err, "undelivered packets") {
		t.Error("0xEB8B96", "wrong error:", err)
	}
} //                                                           Test_Sender_Send_

// (sd *Sender) SendString(name string, s string) error
//
// go test -run Test_Sender_SendString_
//
func Test_Sender_SendString_(t *testing.T) {
	sd := Sender{
		Address:   "127.0.0.0",
		Port:      9876,
		CryptoKey: []byte("12345678901234567890123456789012"),
		Config: &Configuration{
			Cipher:            &aesCipher{},
			Compressor:        &zlibCompressor{},
			PacketSizeLimit:   1024,
			PacketPayloadSize: 512,
			VerboseSender:     true,
			SendRetries:       2,
			ReplyTimeout:      500 * time.Millisecond,
			WriteTimeout:      500 * time.Millisecond,
		},
	}
	err := sd.SendString("greeting", "Hello World!")
	if !matchError(err, "undelivered packets") {
		t.Error("0xEE8E8D", "wrong error:", err)
	}
} //                                                     Test_Sender_SendString_

// -----------------------------------------------------------------------------
// # Informatory Properties (sd *Sender)

// (sd *Sender) AverageResponseMs() float64
//
// go test -run Test_Sender_AverageResponseMs_
//
func Test_Sender_AverageResponseMs_(t *testing.T) {
	var sd Sender
	if n := sd.AverageResponseMs(); n < 0 || n > 0 {
		t.Error("0xE29B40")
	}
	sd.stats.packetsDelivered = 10
	sd.stats.transferTime = time.Millisecond
	if n := sd.AverageResponseMs(); n < 0.1 || n > 0.1 {
		t.Error("0xE01B5F")
	}
} //                                              Test_Sender_AverageResponseMs_

// (sd *Sender) TransferSpeedKBpS() float64
//
// go test -run Test_Sender_TransferSpeedKBpS_
//
func Test_Sender_TransferSpeedKBpS_(t *testing.T) {
	var sd Sender
	if n := sd.TransferSpeedKBpS(); n < 0 || n > 0 {
		t.Error("0xEE99D3")
	}
	sd.stats.transferTime = time.Second
	sd.stats.bytesDelivered = 88 * 1024
	if n := sd.TransferSpeedKBpS(); n < 88 || n > 88 {
		t.Error("0xED0BD8")
	}
} //                                              Test_Sender_TransferSpeedKBpS_

// -----------------------------------------------------------------------------
// # Informatory Methods (sd *Sender)

// (sd *Sender) LogStats(logFunc ...interface{})
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
		sd := Sender{Config: NewDefaultConfig()}
		sd.Config.LogFunc = logPrintln
		sd.packets = []senderPacket{{
			sentHash:      []byte{0},
			confirmedHash: []byte{0},
		}}
		sd.LogStats(logFunc)
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
			t.Error("0xE63A81" + "\n" +
				"expect:\n" + expect + "\n" + "got:\n" + got)
			fmt.Println([]byte(expect))
			fmt.Println([]byte(got))
		}
	}
	test(nil)
	test(logPrintln)
	test(fmtPrintln)
} //                                                       Test_Sender_LogStats_

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (sd *Sender)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) connect() (*net.UDPConn, error)
//
// go test -run Test_Sender_connect_*

// must succeed
func Test_Sender_connect_1(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig(), Address: "127.0.0.1", Port: 9876}
	netDialUDP := func(network string, laddr, raddr *net.UDPAddr,
	) (netUDPConn, error) {
		return net.DialUDP(network, laddr, raddr)
	}
	udpConn, err := sd.connectDI(netDialUDP)
	if udpConn == nil {
		t.Error("0xEC4B85")
	}
	if err != nil {
		t.Error("0xEF79EB", err)
	}
} //                                                       Test_Sender_connect_1

// must fail because Address is invalid
func Test_Sender_connect_2(t *testing.T) {
	sd := Sender{Address: "257.258.259.260", Port: 9876}
	udpConn, err := sd.connect()
	if udpConn != nil {
		t.Error("0xEF8F6B")
	}
	if !matchError(err, "ResolveUDPAddr:") {
		t.Error("0xEC2E79")
	}
} //                                                       Test_Sender_connect_2

// must fail because Port is invalid
func Test_Sender_connect_3(t *testing.T) {
	sd := Sender{Address: "127.0.0.1", Port: 65536}
	udpConn, err := sd.connect()
	if udpConn != nil {
		t.Error("0xE6FA25")
	}
	if !matchError(err, "invalid port") {
		t.Error("0xEA15E2", "wrong error:", err)
	}
} //                                                       Test_Sender_connect_3

// must fail when net.DialUDP() fails
func Test_Sender_connect_4(t *testing.T) {
	sd := Sender{Address: "127.0.0.1", Port: 9876}
	netDialUDP := func(_ string, _, _ *net.UDPAddr) (netUDPConn, error) {
		return nil, makeError(0xEC10B4, "failed netDialUDP")
	}
	udpConn, err := sd.connectDI(netDialUDP)
	if udpConn != nil {
		t.Error("0xED16AA")
	}
	if !matchError(err, "failed netDialUDP") {
		t.Error("0xE2FE6C", "wrong error:", err)
	}
} //                                                       Test_Sender_connect_4

// must fail when conn.SetWriteBuffer() fails
func Test_Sender_connect_5(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig(), Address: "127.0.0.1", Port: 9876}
	netDialUDP := func(_ string, _, _ *net.UDPAddr) (netUDPConn, error) {
		return &mockNetUDPConn{failSetWriteBuffer: true}, nil
	}
	udpConn, err := sd.connectDI(netDialUDP)
	if udpConn != nil {
		t.Error("0xEC77F9")
	}
	if !matchError(err, "failed SetWriteBuffer") {
		t.Error("0xED84F9", "wrong error:", err)
	}
} //                                                       Test_Sender_connect_5

// (sd *Sender) close() error
//
// go test -run Test_Sender_close_
//
func Test_Sender_close_(t *testing.T) {
	var sd Sender
	err := sd.close()
	if err != nil {
		t.Error("0xEE7E05", err)
	}
	sd.conn = makeTestConn()
	err = sd.conn.Close()
	if err != nil {
		t.Error("0xE3DD56", err)
	}
	err = sd.conn.Close()
	if !matchError(err, "use of closed network connection") {
		t.Error("0xE5FE16", "wrong error:", err)
	}
} //                                                          Test_Sender_close_

// -----------------------------------------------------------------------------
// # Internal Helper Methods (sd *Sender)

// (sd *Sender) getPacketCount(length int) int
//
// go test -run Test_Sender_getPacketCount_
//
func Test_Sender_getPacketCount_(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig()}
	sd.Config.PacketPayloadSize = 1000
	//
	if sd.getPacketCount(0) != 0 {
		t.Error("0xE6C4D4")
	}
	if sd.getPacketCount(-100000) != 0 {
		t.Error("0xE81D08")
	}
	if sd.getPacketCount(1) != 1 {
		t.Error("0xE55EB9")
	}
	if sd.getPacketCount(1000) != 1 {
		t.Error("0xE87CB1")
	}
	if sd.getPacketCount(1001) != 2 {
		t.Error("0xE25DD0")
	}
	if sd.getPacketCount(10000) != 10 {
		t.Error("0xEE5EF4")
	}
} //                                                 Test_Sender_getPacketCount_

// (sd *Sender) logError(id uint32, a ...interface{}) error
//
// go test -run Test_Sender_logError_
//
func Test_Sender_logError_(t *testing.T) {
	var msg string
	sd := Sender{Config: NewDefaultConfig()}
	sd.Config.LogFunc = func(a ...interface{}) {
		msg = fmt.Sprintln(a...)
	}
	sd.logError(0xE12345, "abc", 123)
	if msg != "ERROR 0xE12345: abc 123\n" {
		t.Error("0xE5CB5D")
	}
} //                                                       Test_Sender_logError_

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) makePacket(data []byte) (*senderPacket, error)
//
// go test -run Test_Sender_makePacket_*

// must succeed creating a packet to send
func Test_Sender_makePacket_1(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig()}
	packet, err := sd.makePacket([]byte{1, 2, 3})
	if err != nil {
		t.Error("0xE0AE90", err)
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
} //                                                    Test_Sender_makePacket_1

// must fail to create a packet larger than Config.PacketSizeLimit
func Test_Sender_makePacket_2(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig()}
	data := make([]byte, sd.Config.PacketSizeLimit+1)
	packet, err := sd.makePacket(data)
	if packet != nil {
		t.Error("0xE0FE30")
	}
	if !matchError(err, "PacketSizeLimit") {
		t.Error("0xE69EA5", "wrong error:", err)
	}
} //                                                    Test_Sender_makePacket_2

// end
