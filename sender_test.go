// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[sender_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_Sender_*

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// func Send(addr string, k string, v []byte, cryptoKey []byte,
//     config ...*Configuration,
// ) error
//
// go test -run Test_Send_*

// must fail because there are too many 'config' arguments
func Test_Send_A(t *testing.T) {
	cryptoKey := []byte("3z5EdC485Ex9Wy0AsY4Apu6930Bx57Z0")
	var cf *Configuration
	//
	//                addr              k      v        cryptoKey  config
	err := SendString("127.0.0.1:9876", "msg", "test!", cryptoKey, cf, cf)
	if !matchError(err, "too many 'config' arguments") {
		t.Error("0xEE06B6", err)
	}
}

// -----------------------------------------------------------------------------

// SendString(addr string, k, v string, cryptoKey []byte,
//     config ...*Configuration,
// ) error
//
// go test -run Test_SendString_
//
func Test_SendString_(t *testing.T) {
	//
	cryptoKey := []byte("3z5EdC485Ex9Wy0AsY4Apu6930Bx57Z0")
	//
	// set-up and run the receiver
	received := map[string][]byte{} // collects received keys and values
	_, rc := makeConfigAndReceiver(cryptoKey, &received)
	go func() { _ = rc.Run() }()
	defer func() { rc.Stop() }()
	time.Sleep(200 * time.Millisecond)
	//
	err := SendString("127.0.0.1:9876", "_k_", "_v_", cryptoKey, nil)
	if err != nil {
		t.Error("0xE4A1ED", err)
	}
	if len(received) != 1 {
		t.Error("0xED5B82", err)
	}
	for k, v := range received {
		if k != "_k_" || string(v) != "_v_" {
			t.Error("0xEE56FE", err)
		}
	}
}

// -----------------------------------------------------------------------------
// # Main Methods (sd *Sender)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) Send(k string, v []byte) error
//
// go test -run Test_Sender_Send_*

func Test_Sender_Send_1(t *testing.T) {
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
	if !matchError(err, "invalid Sender.Config") {
		t.Error("0xE08E7C", "wrong error:", err)
	}
	sd.Config.PacketSizeLimit = 65535 - 8
	//
	sd.Address = ""
	err = sd.Send("", nil)
	if !matchError(err, "missing Sender.Address") {
		t.Error("0xEC20C3", "wrong error:", err)
	}
	//
	sd.Address = "127.0.0.0:0"
	err = sd.Send("", nil)
	if !matchError(err, "invalid port in Sender.Address") {
		t.Error("0xE24E74", "wrong error:", err)
	}
	sd.Address = "127.0.0.0:9876"
	//
	sd.Config.VerboseSender = true
	sd.Config.SendRetries = 2
	sd.Config.ReplyTimeout = 500 * time.Millisecond
	sd.Config.WriteTimeout = 500 * time.Millisecond
	err = sd.Send("", nil)
	if !matchError(err, "undelivered packets") {
		t.Error("0xEB8B96", "wrong error:", err)
	}
}

func Test_Sender_Send_2(t *testing.T) {
	connect := func() (netUDPConn, error) {
		return nil, makeError(0xEF2DC4, "failed connect")
	}
	sd := makeTestSender()
	err := sd.sendDI("greeting", []byte("Hello!"),
		connect, sd.sendUndeliveredPackets)
	if !matchError(err, "failed connect") {
		t.Error("0xEA93AF")
	}
}

func Test_Sender_Send_3(t *testing.T) {
	sendUndeliveredPackets := func() error {
		return makeError(0xE9AF68, "failed sendUndeliveredPackets")
	}
	sd := makeTestSender()
	err := sd.sendDI("greeting", []byte("Hello!"),
		sd.connect, sendUndeliveredPackets)
	if !matchError(err, "failed sendUndeliveredPackets") {
		t.Error("0xED9E31")
	}
}

// -----------------------------------------------------------------------------

// (sd *Sender) SendString(k, v string) error
//
// go test -run Test_Sender_SendString_
//
func Test_Sender_SendString_(t *testing.T) {
	sd := makeTestSender()
	err := sd.SendString("greeting", "Hello World!")
	if !matchError(err, "undelivered packets") {
		t.Error("0xEE8E8D", "wrong error:", err)
	}
}

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
}

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
}

// -----------------------------------------------------------------------------
// # Informatory Methods (sd *Sender)

// (sd *Sender) LogStats(w io.Writer)
//
// go test -run Test_Sender_LogStats_
//
func Test_Sender_LogStats_(t *testing.T) {
	var tlog strings.Builder
	//
	sd := Sender{Config: NewDefaultConfig()}
	sd.Config.LogWriter = &tlog
	sd.packets = []senderPacket{
		{sentHash: []byte{0x0}, confirmedHash: []byte{0x0}},
		{sentHash: []byte{0x1}, confirmedHash: []byte{0x0}},
	}
	sd.stats.bytesDelivered = 123000
	sd.stats.bytesLost = 456
	sd.stats.packetsDelivered = 10
	sd.stats.packetsLost = 1
	sd.stats.transferTime = 300 * time.Millisecond
	// -----------------
	sd.LogStats(&tlog)
	// -----------------
	want := "" +
		"SN: 0    T0: 0001-01-01 00:00:00 +000 T1: NONE ✔ 0.0 ms\n" +
		"SN: 1    T0: 0001-01-01 00:00:00 +000 T1: NONE LOST 0.0 ms\n" +
		"B. delivered: 123000\n" +
		"Bytes lost  : 456\n" +
		"P. delivered: 10\n" +
		"Packets lost: 1\n" +
		"Time in item: 0.3 s\n" +
		"Avg./ Packet: 30.0 ms\n" +
		"Trans. speed: 400.0 KiB/s\n"
	//
	got := tlog.String()
	if got != want {
		t.Error("0xE63A81\n" + "want:\n" + want + "\n" + "got:\n" + got)
		// fmt.Println("want bytes:", []byte(want))
		// fmt.Println(" got bytes:", []byte(got))
	}
}

// -----------------------------------------------------------------------------
// # Internal Lifecycle Methods (sd *Sender)

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) connect() (netUDPConn, error)
//
// go test -run Test_Sender_connect_*

// must succeed
func Test_Sender_connect_1(t *testing.T) {
	sd := Sender{Config: NewDefaultConfig(), Address: "127.0.0.1:9876"}
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
}

// must fail because the host in the address is invalid
func Test_Sender_connect_2(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "257.258.259.260:9876"
	udpConn, err := sd.connect()
	if udpConn != nil {
		t.Error("0xEF8F6B")
	}
	if !matchError(err, "ResolveUDPAddr:") {
		t.Error("0xEC2E79")
	}
}

// must fail because the port in the address is invalid
func Test_Sender_connect_3(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "127.0.0.1:65536"
	udpConn, err := sd.connect()
	if udpConn != nil {
		t.Error("0xE6FA25")
	}
	if !matchError(err, "invalid port") {
		t.Error("0xEA15E2", "wrong error:", err)
	}
}

// must fail when net.DialUDP() fails
func Test_Sender_connect_4(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "127.0.0.1:9876"
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
}

// must fail when conn.SetWriteBuffer() fails
func Test_Sender_connect_5(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "127.0.0.1:9876"
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
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) close() error
//
// go test -run Test_Sender_close_*

// must succeed
func Test_Sender_close_1(t *testing.T) {
	var tlog strings.Builder
	sd := makeTestSender()
	sd.Config = NewDebugConfig(&tlog)
	sd.close()
	if sd.conn != nil {
		t.Error("0xE1FE10")
	}
	ts := tlog.String()
	if ts != "" {
		t.Error("0xE0AA32")
	}
	sd.close()
	if sd.conn != nil {
		t.Error("0xE67DB6")
	}
	ts = tlog.String()
	if ts != "" {
		t.Error("0xEA7A80")
	}
}

// must write to log when sd.conn.Close() fails
func Test_Sender_close_2(t *testing.T) {
	var tlog strings.Builder
	sd := makeTestSender()
	sd.conn = &mockNetUDPConn{failClose: true}
	sd.Config = NewDebugConfig(&tlog)
	sd.close()
	ts := tlog.String()
	if !strings.Contains(ts, "failed Close") {
		t.Error("0xEA8D88")
	}
}

// -----------------------------------------------------------------------------
// # Internal Helper Methods (sd *Sender)

// (sd *Sender) logError(id uint32, a ...interface{}) error
//
// go test -run Test_Sender_logError_
//
func Test_Sender_logError_(t *testing.T) {
	var tlog strings.Builder
	sd := makeTestSender()
	sd.Config.LogWriter = &tlog
	sd.logError(0xE12345, "abc", 123)
	//
	ts := tlog.String()
	if ts != "ERROR 0xE12345: abc 123\n" {
		t.Error("0xE5CB5D")
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) makePacket(data []byte) (*senderPacket, error)
//
// go test -run Test_Sender_makePacket_*

// must succeed creating a packet to send
func Test_Sender_makePacket_1(t *testing.T) {
	sd := makeTestSender()
	pk, err := sd.makePacket([]byte{1, 2, 3})
	if err != nil {
		t.Error("0xE0AE90", err)
	}
	wantData := []byte{1, 2, 3}
	if !bytes.Equal(pk.data, wantData) {
		t.Error("0xEF4E82")
	}
	wantSentHash := getHash([]byte{1, 2, 3})
	if !bytes.Equal(pk.sentHash, wantSentHash) {
		t.Error("0xE51B95")
	}
	n := time.Since(pk.sentTime)
	if n > time.Millisecond {
		t.Error("0xE1FA4B")
	}
	if pk.confirmedHash != nil {
		t.Error("0xEA0E4B")
	}
	if !pk.confirmedTime.IsZero() {
		t.Error("0xE21EB4")
	}
}

// must fail to create a packet larger than Config.PacketSizeLimit
func Test_Sender_makePacket_2(t *testing.T) {
	sd := makeTestSender()
	data := make([]byte, sd.Config.PacketSizeLimit+1)
	pk, err := sd.makePacket(data)
	if pk != nil {
		t.Error("0xE0FE30")
	}
	if !matchError(err, "PacketSizeLimit") {
		t.Error("0xE69EA5", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (sd *Sender) validateAddress() error
//
// go test -run Test_Sender_validateAddress_*

// must return nil when Address is valid
func Test_Sender_validateAddress_1(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "127.0.0.1:9876"
	err := sd.validateAddress()
	if err != nil {
		t.Error("0xEC9A5E")
	}
}

// must return error "missing Sender.Address" when Address is ""
func Test_Sender_validateAddress_2(t *testing.T) {
	sd := makeTestSender()
	sd.Address = ""
	err := sd.validateAddress()
	if !matchError(err, "missing Sender.Address") {
		t.Error("0xEF8A89")
	}
}

// must return error "missing Sender.Address" when Address is "\r \n"
func Test_Sender_validateAddress_3(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "\r \n"
	err := sd.validateAddress()
	if !matchError(err, "missing Sender.Address") {
		t.Error("0xEB66F3")
	}
}

// must return error "missing Sender.Address" when port is not specified
func Test_Sender_validateAddress_4(t *testing.T) {
	sd := makeTestSender()
	sd.Address = "127.0.0.1"
	err := sd.validateAddress()
	if !matchError(err, "invalid port in Sender.Address") {
		t.Error("0xE45F34")
	}
}

// -----------------------------------------------------------------------------

// makeConfigAndReceiver creates and returns a
// Configuration and Receiver for testing Sender.
func makeConfigAndReceiver(cryptoKey []byte, received *map[string][]byte,
) (*Configuration, *Receiver) {
	cf := NewDefaultConfig()
	cf.ReplyTimeout = 250 * time.Millisecond
	cf.WriteTimeout = 250 * time.Millisecond
	//
	rc := Receiver{Port: 9876, CryptoKey: cryptoKey, Config: cf,
		Receive: func(k string, v []byte) error {
			(*received)[k] = []byte(v)
			return nil
		},
	}
	return cf, &rc
}

// makeTestSender creates a properly-configured Sender for testing.
func makeTestSender() *Sender {
	cf := Configuration{
		Cipher:            &aesCipher{},
		Compressor:        &zlibCompressor{},
		PacketSizeLimit:   1024,
		PacketPayloadSize: 512,
		VerboseSender:     true,
		SendRetries:       2,
		ReplyTimeout:      500 * time.Millisecond,
		WriteTimeout:      500 * time.Millisecond,
	}
	sd := Sender{
		Address:   "127.0.0.0:9876",
		CryptoKey: []byte("12345678901234567890123456789012"),
		Config:    &cf,
	}
	cf.Cipher.SetKey(sd.CryptoKey)
	return &sd
}

// end
