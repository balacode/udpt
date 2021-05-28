// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[receiver_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_Receiver_*

// -----------------------------------------------------------------------------

// newRunnableReceiver() creates a Receiver with all required fields set
func newRunnableReceiver() Receiver {
	ret := Receiver{
		Port:        9876,
		CryptoKey:   []byte("0123456789abcdefghijklmnopqrst12"),
		Config:      NewDefaultConfig(),
		ReceiveData: func(k string, v []byte) error { return nil },
		ProvideData: func(k string) ([]byte, error) { return nil, nil },
	}
	ret.Config.ReplyTimeout = 500 * time.Millisecond
	ret.Config.WriteTimeout = 500 * time.Millisecond
	return ret
}

// -----------------------------------------------------------------------------
// (rc *Receiver) Run() error
//
// go test -run Test_Receiver_Run_*

// expecting no startup error
func Test_Receiver_Run_01(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Config = nil
	go func() {
		time.Sleep(500 * time.Millisecond)
		rc.Stop()
	}()
	err := rc.Run()
	if rc.Config == nil {
		t.Error("0xEF66E0")
	}
	if err != nil {
		t.Error("0xEF9D95", err)
	}
}

// must fail to start: Config.Cipher is not specfied
func Test_Receiver_Run_02(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Config.Cipher = nil
	err := rc.Run()
	if !matchError(err, "nil Configuration.Cipher") {
		t.Error("0xE4F1AF", "wrong error:", err)
	}
}

// must fail to start: Config is invalid
func Test_Receiver_Run_03(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Config.PacketSizeLimit = 0
	err := rc.Run()
	if !matchError(err, "invalid Configuration.PacketSizeLimit") {
		t.Error("0xEC1D86", "wrong error:", err)
	}
}

// must fail to start: CryptoKey is not specfied
func Test_Receiver_Run_04(t *testing.T) {
	rc := newRunnableReceiver()
	rc.CryptoKey = nil
	err := rc.Run()
	if !matchError(err, "Receiver.CryptoKey") {
		t.Error("0xE57F1E", "wrong error:", err)
	}
}

// must fail to start: CryptoKey is wrong size
func Test_Receiver_Run_05(t *testing.T) {
	rc := newRunnableReceiver()
	rc.CryptoKey = []byte{1, 2, 3}
	err := rc.Run()
	if !matchError(err, "AES-256 key") {
		t.Error("0xE19A88", "wrong error:", err)
	}
}

// must fail to start: Port is not set
func Test_Receiver_Run_06(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Port = 0
	err := rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xE21D17", "wrong error:", err)
	}
}

// must fail to start: Port number is too high
func Test_Receiver_Run_07(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Port = 65535 + 1
	err := rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xE8E6D5", "wrong error:", err)
	}
}

// must fail to start: Port number is negative
func Test_Receiver_Run_08(t *testing.T) {
	rc := newRunnableReceiver()
	rc.Port = -123
	err := rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xED3AE1", "wrong error:", err)
	}
}

// must fail to start: ReceiveData is not specified
func Test_Receiver_Run_09(t *testing.T) {
	rc := newRunnableReceiver()
	rc.ReceiveData = nil
	err := rc.Run()
	if !matchError(err, "nil Receiver.ReceiveData") {
		t.Error("0xE7C0AC", "wrong error:", err)
	}
}

// must fail to start: ProvideData is not specified
func Test_Receiver_Run_10(t *testing.T) {
	rc := newRunnableReceiver()
	rc.ProvideData = nil
	err := rc.Run()
	if !matchError(err, "nil Receiver.ProvideData") {
		t.Error("0xEF5FF2", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) Stop()
//
// go test -run Test_Receiver_Stop_*

func Test_Receiver_Stop_1(t *testing.T) {
	var rc Receiver
	rc.Stop()
	if rc.conn != nil {
		t.Error("0xE1FC2F")
	}
}

func Test_Receiver_Stop_2(t *testing.T) {
	var err error
	fn := func(a ...interface{}) {
		err = makeError(0xEC8A02, a...)
	}
	rc := Receiver{
		Config: NewDebugConfig(fn),
		conn:   &mockNetUDPConn{failClose: true},
	}
	rc.Stop()
	if rc.conn != nil {
		t.Error("0xE5D85F")
	}
	if !matchError(err, "failed Close") {
		t.Error("0xEF17A7", "wrong error:", err)
	}
}

// -----------------------------------------------------------------------------
// # Run() Helpers

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) initRun() error

// must succeed
func Test_Receiver_initRun_1(t *testing.T) {
	netResolveUDPAddr := func(string, string) (*net.UDPAddr, error) {
		return nil, nil
	}
	netListenUDP := func(string, *net.UDPAddr) (*net.UDPConn, error) {
		return &net.UDPConn{}, nil
	}
	rc := newRunnableReceiver()
	cons := ""
	rc.Config.VerboseReceiver = true
	rc.Config.LogFunc = func(a ...interface{}) {
		cons += fmt.Sprintln(a...)
	}
	err := rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if rc.Config == nil {
		t.Error("0xE9A35E")
	}
	if err != nil {
		t.Error("0xE0A70B", err)
	}
	if !reflect.DeepEqual(rc.conn, &net.UDPConn{}) {
		t.Error("0xEC5C19")
	}
	if !strings.Contains(cons, strings.Repeat("-", 80)) {
		t.Error("0xEA76F8")
	}
	if !strings.Contains(cons, "Receiver listening...") {
		t.Error("0xEE15D5")
	}
}

// must succeed
func Test_Receiver_initRun_2(t *testing.T) {
	var c1, c2 bool
	netResolveUDPAddr :=
		func(network string, addr string) (*net.UDPAddr, error) {
			c1 = true
			if network != "udp" {
				t.Error("0xE94D6D")
			}
			if addr != "0.0.0.0:9876" {
				t.Error("0xED87E1")
			}
			return &net.UDPAddr{IP: []byte{5, 4, 3, 2}, Port: 9876}, nil
		}
	netListenUDP :=
		func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
			c2 = true
			if network != "udp" {
				t.Error("0xEE3B32")
			}
			if !reflect.DeepEqual(
				laddr, &net.UDPAddr{IP: []byte{5, 4, 3, 2}, Port: 9876},
			) {
				t.Error("0xE70A3D")
			}
			return &net.UDPConn{}, nil
		}
	rc := newRunnableReceiver()
	cons := ""
	rc.Config.VerboseReceiver = true
	rc.Config.LogFunc = func(a ...interface{}) {
		cons += fmt.Sprintln(a...)
	}
	rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !c1 {
		t.Error("0xE31BE4")
	}
	if !c2 {
		t.Error("0xE73C32")
	}
}

// fail because Config is not set
func Test_Receiver_initRun_3(t *testing.T) {
	rc := Receiver{Config: &Configuration{}}
	err := rc.initRun()
	if err == nil {
		t.Error("0xEA74DE", "err must not be nil")
	}
}

// fail because Config is not valid
func Test_Receiver_initRun_4(t *testing.T) {
	var c1, c2 bool
	netResolveUDPAddr := func(string, string) (*net.UDPAddr, error) {
		c1 = true
		return nil, nil
	}
	netListenUDP := func(string, *net.UDPAddr) (*net.UDPConn, error) {
		c2 = true
		return nil, nil
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.PacketSizeLimit = -1
	err := rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "invalid Configuration.PacketSizeLimit") {
		t.Error("0xED9BC3", "wrong error:", err)
	}
	// at this point, none of the net functions must have been called
	if c1 {
		t.Error("0xE22A08")
	}
	if c2 {
		t.Error("0xE26B1E")
	}
}

// must fail because Receiver.Port is wrong
func Test_Receiver_initRun_5(t *testing.T) {
	var c1, c2 bool
	netResolveUDPAddr := func(string, string) (*net.UDPAddr, error) {
		c1 = true
		return nil, nil
	}
	netListenUDP := func(string, *net.UDPAddr) (*net.UDPConn, error) {
		c2 = true
		return nil, nil
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.PacketSizeLimit = 2048
	//
	// Port is not set
	err := rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "Receiver.Port") {
		t.Error("0xEE7BF2", "wrong error:", err)
	}
	// Port is out of range
	rc.Port = -789
	if !matchError(err, "Receiver.Port") {
		t.Error("0xEF50CF", "wrong error:", err)
	}
	rc.Port = 65536
	if !matchError(err, "Receiver.Port") {
		t.Error("0xEE72E1", "wrong error:", err)
	}
	// at this point, none of the net functions must have been called
	if c1 {
		t.Error("0xEE4CD2")
	}
	if c2 {
		t.Error("0xEA0C85")
	}
}

// must fail because CryptoKey is wrong
func Test_Receiver_initRun_6(t *testing.T) {
	var c1, c2 bool
	netResolveUDPAddr := func(string, string) (*net.UDPAddr, error) {
		c1 = true
		return nil, nil
	}
	netListenUDP := func(string, *net.UDPAddr) (*net.UDPConn, error) {
		c2 = true
		return nil, nil
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Port = 9876
	//
	// CryptoKey is not set
	err := rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "Receiver.CryptoKey") {
		t.Error("0xEF28D6", "wrong error:", err)
	}
	// CryptoKey is too short
	rc.CryptoKey = []byte{1, 2, 3}
	err = rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "Receiver.CryptoKey") {
		t.Error("0xED60C9", "wrong error:", err)
	}
	// CryptoKey is too long
	rc.CryptoKey = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
		18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33,
	}
	err = rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "Receiver.CryptoKey") {
		t.Error("0xE6BF78", "wrong error:", err)
	}
	// at this point, none of the net functions must have been called
	if c1 {
		t.Error("0xE78C71")
	}
	if c2 {
		t.Error("0xE6AC5D")
	}
}

// must fail because ReceiveData or ProvideData is not assigned
func Test_Receiver_initRun_7(t *testing.T) {
	var c1, c2 bool
	netResolveUDPAddr := func(string, string) (*net.UDPAddr, error) {
		c1 = true
		return nil, nil
	}
	netListenUDP := func(string, *net.UDPAddr) (*net.UDPConn, error) {
		c2 = true
		return nil, nil
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Port = 9876
	rc.CryptoKey = []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17,
		18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
	}
	// fail because ReceiveData is not assigned
	err := rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "nil Receiver.ReceiveData") {
		t.Error("0xE2F00D", "wrong error:", err)
	}
	rc.ReceiveData = func(k string, v []byte) error { return nil }
	//
	// fail because ProvideData is not assigned
	err = rc.initRunDI(netResolveUDPAddr, netListenUDP)
	if !matchError(err, "nil Receiver.ProvideData") {
		t.Error("0xE3E29F", "wrong error:", err)
	}
	rc.ProvideData = func(k string) ([]byte, error) { return nil, nil }
	//
	// at this point, none of the net functions must have been called
	if c1 {
		t.Error("0xE8F58D")
	}
	if c2 {
		t.Error("0xE1B17E")
	}
}

// must fail when netResolveUDPAddr fails
func Test_Receiver_initRun_8(t *testing.T) {
	netResolveUDPAddr :=
		func(network string, addr string) (*net.UDPAddr, error) {
			return nil, makeError(0xE2F60A, "failed netResolveUDPAddr")
		}
	rc := newRunnableReceiver()
	err := rc.initRunDI(netResolveUDPAddr, net.ListenUDP)
	if rc.Config == nil {
		t.Error("0xE8B44D")
	}
	if !matchError(err, "failed netResolveUDPAddr") {
		t.Error("0xE44EF9", "wrong error:", err)
	}
}

// must fail when netListenUDP fails
func Test_Receiver_initRun_9(t *testing.T) {
	netListenUDP :=
		func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
			return nil, makeError(0xE2F33D, "failed netListenUDP")
		}
	rc := newRunnableReceiver()
	err := rc.initRunDI(net.ResolveUDPAddr, netListenUDP)
	if rc.Config == nil {
		t.Error("0xE51C18")
	}
	if !matchError(err, "failed netListenUDP") {
		t.Error("0xE9D64F", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) buildReply(recv []byte) (reply []byte, err error)

// must succeed
func Test_Receiver_buildReply_1(t *testing.T) {
	zc := &zlibCompressor{}
	comp, err := zc.Compress([]byte("abc"))
	if err != nil {
		t.Error("0xE29CE8", err)
	}
	rc := Receiver{Config: NewDefaultConfig()}
	logMsg, recKey, recVal := "", "", ""
	rc.Config.Cipher.SetKey([]byte(testAESKey))
	rc.Config.LogFunc = func(a ...interface{}) {
		logMsg += fmt.Sprintln(a...)
	}
	rc.ReceiveData = func(k string, v []byte) error {
		recKey, recVal = k, string(v)
		return nil
	}
	reply, err := rc.buildReply([]byte(
		tagFragment + "key:test1 " +
			"hash:BA7816BF8F01CFEA414140DE5DAE2223" +
			"B00361A396177A9CB410FF61F20015AD sn:1 count:1\n" +
			string(comp),
	))
	if recKey != "test1" {
		t.Error("0xE89AA5")
	}
	if recVal != "abc" {
		t.Error("0xE07D6B")
	}
	if len(reply) < 10 {
		t.Error("0xE5AB18")
	}
	if err != nil {
		t.Error("0xE06F48")
	}
	if !strings.Contains(logMsg, "received: test1") {
		t.Error("0xE70F40", "wrong reply:", logMsg)
	}
}

// must fail because sent data is nil
func Test_Receiver_buildReply_2(t *testing.T) {
	logErrorMsg := ""
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.LogFunc = func(a ...interface{}) {
		logErrorMsg += fmt.Sprintln(a...)
	}
	reply, err := rc.buildReply(nil)
	if reply != nil {
		t.Error("0xE18DB7")
	}
	if err != nil {
		t.Error("0xEC58A7")
	}
	if !strings.Contains(logErrorMsg, "received no data") {
		t.Error("0xE75A71", "wrong error:", logErrorMsg)
	}
}

// must fail because packet header is invalid
func Test_Receiver_buildReply_3(t *testing.T) {
	logErrorMsg := ""
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.Cipher.SetKey([]byte(testAESKey))
	rc.Config.LogFunc = func(a ...interface{}) {
		logErrorMsg += fmt.Sprintln(a...)
	}
	reply, err := rc.buildReply([]byte("XYZ: ..."))
	if string(reply) != "invalid_packet_header" {
		t.Error("0xE2CA90")
	}
	if !matchError(err, "invalid packet header") {
		t.Error("0xEC3D21", "wrong error:", err)
	}
	if !strings.Contains(logErrorMsg, "invalid packet header") {
		t.Error("0xE53F59", "wrong error:", logErrorMsg)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) sendReply(conn netUDPConn, addr net.Addr, reply []byte)

func Test_Receiver_sendReply_1(t *testing.T) {
	var (
		rc   = Receiver{Config: NewDefaultConfig()}
		cn   = &mockNetUDPConn{}
		addr = &net.UDPAddr{IP: []byte{127, 0, 0, 0}, Port: 9876}
		cons = ""
	)
	rc.Config.VerboseReceiver = true
	rc.Config.WriteTimeout = 7 * time.Millisecond
	rc.Config.LogFunc = func(a ...interface{}) {
		cons += fmt.Sprintln(a...)
	}
	deadline := time.Now().Add(7 * time.Millisecond)
	//
	rc.sendReply(cn, addr, []byte{1, 2, 3, 4, 5})
	//
	since := cn.writeDeadline.Sub(deadline)
	if since > time.Millisecond || cn.writeDeadline.IsZero() {
		t.Error("0xE2C11E")
	}
	rc.sendReply(cn, addr, []byte{6, 7, 8, 9, 10})
	//
	if cn.nSetWriteDeadline != 2 {
		t.Error("0xE77AF3")
	}
	if cn.nWriteTo != 2 {
		t.Error("0xEE1E5F")
	}
	if (cn.nSetReadDeadline + cn.nReadFrom + cn.nClose) != 0 {
		t.Error("0xE35FD7")
	}
	if !bytes.Equal(cn.written, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("0xE07FA9")
	}
	if strings.Count(cons, "Receiver wrote 5 bytes to 127.0.0.0:9876") != 2 {
		t.Error("0xEA1AE3")
	}
}

func Test_Receiver_sendReply_2(t *testing.T) {
	var (
		rc   = Receiver{Config: NewDefaultConfig()}
		cn   = &mockNetUDPConn{failSetWriteDeadline: true}
		addr = &net.UDPAddr{IP: []byte{127, 0, 0, 0}, Port: 9876}
	)
	rc.sendReply(cn, addr, nil)
	if cn.nSetWriteDeadline != 1 {
		t.Error("0xED1AD3")
	}
	if (cn.nSetReadDeadline + cn.nReadFrom + cn.nWriteTo + cn.nClose) != 0 {
		t.Error("0xED5E71")
	}
}

func Test_Receiver_sendReply_3(t *testing.T) {
	var (
		rc   = Receiver{Config: NewDefaultConfig()}
		cn   = &mockNetUDPConn{failWriteTo: true}
		addr = &net.UDPAddr{IP: []byte{127, 0, 0, 0}, Port: 9876}
	)
	rc.sendReply(cn, addr, nil)
	if cn.nSetWriteDeadline != 1 {
		t.Error("0xEC75DC")
	}
	if cn.nWriteTo != 1 {
		t.Error("0xEB7CD8")
	}
	if (cn.nSetReadDeadline + cn.nReadFrom + cn.nClose) != 0 {
		t.Error("0xE87A3F")
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) receiveFragment(recv []byte) ([]byte, error)
//
// go test -run Test_Receiver_receiveFragment_*

func Test_Receiver_receiveFragment_01(t *testing.T) {
	var rc Receiver
	data, err := rc.receiveFragment([]byte{})
	if data != nil {
		t.Error("0xE36A92")
	}
	if !matchError(err, "missing header") {
		t.Error("0xEF7AE2", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_02(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.Cipher = nil
	data, err := rc.receiveFragment([]byte(tagFragment))
	if data != nil {
		t.Error("0xE6F51D")
	}
	if !matchError(err, "nil Configuration.Cipher") {
		t.Error("0xE90F36", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_03(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment))
	if data != nil {
		t.Error("0xE9F5CF")
	}
	if !matchError(err, "newline not found") {
		t.Error("0xE8DC8E", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_04(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:0 sn:bad count:1\n"))
	if data != nil {
		t.Error("0xEA0B81")
	}
	if !matchError(err, "bad 'sn'") {
		t.Error("0xEC2C48", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_05(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:0 sn:1 count:bad\n"))
	if data != nil {
		t.Error("0xEA9D01")
	}
	if !matchError(err, "bad 'count'") {
		t.Error("0xEA33B6", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_06(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:0 sn:2 count:1\n"))
	if data != nil {
		t.Error("0xEB21B0")
	}
	if !matchError(err, "out of range") {
		t.Error("0xEB9A96", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_07(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:0 sn:1 count:1\n"))
	if data != nil {
		t.Error("0xE11DF3")
	}
	if !matchError(err, "hex") {
		t.Error("0xE43F0E", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_08(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:GG sn:1 count:1\n"))
	if data != nil {
		t.Error("0xEF09EC")
	}
	if !matchError(err, "hex") {
		t.Error("0xEF4C9B", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_09(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:FF sn:1 count:1\n"))
	if data != nil {
		t.Error("0xE24F86")
	}
	if !matchError(err, "bad hash size") {
		t.Error("0xEB87BE", "wrong error:", err)
	}
}

func Test_Receiver_receiveFragment_10(t *testing.T) {
	rc := Receiver{Config: NewDefaultConfig()}
	data, err := rc.receiveFragment([]byte(tagFragment +
		"key:abc hash:" +
		"12345678123456781234567812345678" +
		"12345678123456781234567812345678 sn:1 count:1\n"))
	if data != nil {
		t.Error("0xE85E88")
	}
	if !matchError(err, "received no data") {
		t.Error("0xE13A6F", "wrong error:", err)
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) sendDataItemHash(req []byte) ([]byte, error)
//
// go test -run Test_Receiver_sendDataItemHash_*

func Test_Receiver_sendDataItemHash_1(t *testing.T) {
	var rc Receiver
	data, err := rc.sendDataItemHash([]byte{})
	if data != nil {
		t.Error("0xE30A2F")
	}
	if !matchError(err, "missing header") {
		t.Error("0xED4B27", "wrong error:", err)
	}
}

func Test_Receiver_sendDataItemHash_2(t *testing.T) {
	var rc Receiver
	data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
	if data != nil {
		t.Error("0xED8A3E")
	}
	if !matchError(err, "nil ProvideData") {
		t.Error("0xE65C25", "wrong error:", err)
	}
}

func Test_Receiver_sendDataItemHash_3(t *testing.T) {
	var rc Receiver
	rc.ProvideData = func(k string) ([]byte, error) {
		return nil, errors.New("test error")
	}
	data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
	if data != nil {
		t.Error("0xEA5B15")
	}
	if !matchError(err, "test error") {
		t.Error("0xEE3C84", "wrong error:", err)
	}
}

func Test_Receiver_sendDataItemHash_4(t *testing.T) {
	var rc Receiver
	rc.ProvideData = func(k string) ([]byte, error) {
		return nil, nil
	}
	data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
	if string(data) != "HASH:"+
		"E3B0C44298FC1C149AFBF4C8996FB924"+
		"27AE41E4649B934CA495991B7852B855" {
		t.Error("0xE8F93C")
	}
	if err != nil {
		t.Error("0xE0D7A2", err)
	}
}

func Test_Receiver_sendDataItemHash_5(t *testing.T) {
	var rc Receiver
	rc.ProvideData = func(k string) ([]byte, error) {
		return []byte("0123456789"), nil
	}
	data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
	if string(data) != "HASH:"+
		"84D89877F0D4041EFB6BF91A16F0248F"+
		"2FD573E6AF05C19F96BEDB9F882F7882" {
		t.Error("0xE37BD7")
	}
	if err != nil {
		t.Error("0xEF13C6", err)
	}
}

// -----------------------------------------------------------------------------
// # Logging Methods

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) logError(a ...interface{})
//
// go test -run Test_Receiver_logError_*

func Test_Receiver_logError_1(t *testing.T) {
	var sb strings.Builder
	var rc Receiver
	//
	rc.logError(0xE12345, "error message")
	//
	got := sb.String()
	if got != "" {
		t.Error("0xE94FB3")
	}
}

func Test_Receiver_logError_2(t *testing.T) {
	var sb strings.Builder
	fn := func(a ...interface{}) {
		sb.WriteString(fmt.Sprint(a...))
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.LogFunc = fn
	//
	rc.logError(0xE12345, "error message")
	//
	got := sb.String()
	if got != "ERROR 0xE12345: error message" {
		t.Error("0xE0FA6C")
	}
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// (rc *Receiver) logInfo(a ...interface{})
//
// go test -run Test_Receiver_logInfo_*

func Test_Receiver_logInfo_1(t *testing.T) {
	var sb strings.Builder
	var rc Receiver
	//
	rc.logInfo("info message")
	//
	got := sb.String()
	if got != "" {
		t.Error("0xEF3F1C")
	}
}

func Test_Receiver_logInfo_2(t *testing.T) {
	var sb strings.Builder
	fn := func(a ...interface{}) {
		sb.WriteString(fmt.Sprint(a...))
	}
	rc := Receiver{Config: NewDefaultConfig()}
	rc.Config.LogFunc = fn
	//
	rc.logInfo("info message")
	//
	got := sb.String()
	if got != "info message" {
		t.Error("0xE74A75")
	}
}

// end
