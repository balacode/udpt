// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[receiver_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run _Receiver_

// -----------------------------------------------------------------------------

// (ob *Receiver) Run() error
//
// go test -run _Receiver_Run_
//
func Test_Receiver_Run_(t *testing.T) {
	var rc Receiver
	//
	// expecting startup error: CryptoKey is not specfied
	err := rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.CryptoKey") {
		t.Error("0xE57F1E")
	}
	// expecting startup error: CryptoKey is wrong size
	rc.CryptoKey = []byte{1, 2, 3}
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "AES-256 key") {
		t.Error("0xE19A88")
	}
	// expecting startup error: Port is not set
	rc.CryptoKey = []byte("0123456789abcdefghijklmnopqrst12")
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.Port") {
		t.Error("0xE21D17")
	}
	// expecting startup error: Port number is wrong
	rc.Port = 65535 + 1
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.Port") {
		t.Error("0xE8E6D5")
	}
	// expecting startup error: Port number is wrong
	rc.Port = 0
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.Port") {
		t.Error("0xED3AE1")
	}
	// expecting startup error: ReceiveData not specified
	rc.Port = 9874
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.ReceiveData") {
		t.Error("0xE7C0AC")
	}
	// expecting startup error: ProvideData not specified
	rc.ReceiveData = func(name string, data []byte) error { return nil }
	err = rc.Run()
	if err == nil || !strings.Contains(err.Error(), "Receiver.ProvideData") {
		t.Error("0xEF5FF2")
	}
	// expecting no startup error
	rc.ProvideData = func(name string) ([]byte, error) { return nil, nil }
	go func() {
		time.Sleep(1 * time.Second)
		rc.Stop()
	}()
	err = rc.Run()
	if err != nil {
		t.Error("0xEF9D95")
	}
} //                                                          Test_Receiver_Run_

// -----------------------------------------------------------------------------
// # Logging Methods

// (ob *Receiver) logError(args ...interface{})
//
// go test -run Test_Receiver_logError_
//
func Test_Receiver_logError_(t *testing.T) {
	var sb strings.Builder
	fn := func(args ...interface{}) {
		sb.WriteString(fmt.Sprint(args...))
	}
	{
		var rc Receiver
		rc.logError(0xE12345, "error message")
		got := sb.String()
		if got != "" {
			t.Error("0xE94FB3")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		rc.Config.LogFunc = fn
		rc.logError(0xE12345, "error message")
		got := sb.String()
		if got != "ERROR 0xE12345: error message" {
			t.Error("0xE0FA6C")
		}
	}
} //                                                     Test_Receiver_logError_

// (ob *Receiver) logInfo(args ...interface{})
//
// go test -run Test_Receiver_logInfo_
//
func Test_Receiver_logInfo_(t *testing.T) {
	var sb strings.Builder
	fn := func(args ...interface{}) {
		sb.WriteString(fmt.Sprint(args...))
	}
	{
		var rc Receiver
		rc.logInfo("info message")
		got := sb.String()
		if got != "" {
			t.Error("0xEF3F1C")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		rc.Config.LogFunc = fn
		rc.logInfo("info message")
		got := sb.String()
		if got != "info message" {
			t.Error("0xE74A75")
		}
	}
} //                                                      Test_Receiver_logInfo_

// end
