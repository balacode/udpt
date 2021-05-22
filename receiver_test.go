// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                  /[receiver_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_Receiver_*

// -----------------------------------------------------------------------------

// (rc *Receiver) Run() error
//
// go test -run Test_Receiver_Run_
//
func Test_Receiver_Run_(t *testing.T) {
	var rc Receiver
	//
	// expecting startup error: Port is not specfied
	err := rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xE57F1E")
	}
	// expecting startup error: Port number too big
	rc.Port = 65535 + 1
	err = rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xE8E6D5")
	}
	// expecting startup error: Port number is negative
	rc.Port = -123
	err = rc.Run()
	if !matchError(err, "Receiver.Port") {
		t.Error("0xED3AE1")
	}
	// expecting startup error: CryptoKey is not set
	rc.Port = 9876
	rc.CryptoKey = nil
	err = rc.Run()
	if !matchError(err, "AES-256 key") {
		t.Error("0xE19A88")
	}
	// expecting startup error: CryptoKey is wrong size
	rc.CryptoKey = []byte{1, 2, 3}
	err = rc.Run()
	if !matchError(err, "AES-256 key") {
		t.Error("0xE19A88")
	}
	// expecting startup error: ReceiveData not specified
	rc.CryptoKey = []byte("0123456789abcdefghijklmnopqrst12")
	rc.Port = 9874
	err = rc.Run()
	if !matchError(err, "nil Receiver.ReceiveData") {
		t.Error("0xE7C0AC")
	}
	// expecting startup error: ProvideData not specified
	rc.ReceiveData = func(name string, data []byte) error { return nil }
	err = rc.Run()
	if !matchError(err, "nil Receiver.ProvideData") {
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

// (rc *Receiver) receiveFragment(recv []byte) ([]byte, error)
//
// go test -run Test_Receiver_receiveFragment_
//
func Test_Receiver_receiveFragment_(t *testing.T) {
	{
		var rc Receiver
		data, err := rc.receiveFragment([]byte{})
		if data != nil {
			t.Error("0xE36A92")
		}
		if !matchError(err, "missing header") {
			t.Error("0xEF7AE2")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		rc.Config.Cipher = nil
		data, err := rc.receiveFragment([]byte(tagFragment))
		if data != nil {
			t.Error("0xE6F51D")
		}
		if !matchError(err, "nil Configuration.Cipher") {
			t.Error("0xE90F36")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment))
		if data != nil {
			t.Error("0xE9F5CF")
		}
		if !matchError(err, "newline not found") {
			t.Error("0xE8DC8E")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:0 sn:bad count:1\n"))
		if data != nil {
			t.Error("0xEA0B81")
		}
		if !matchError(err, "bad 'sn'") {
			t.Error("0xEC2C48")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:0 sn:1 count:bad\n"))
		if data != nil {
			t.Error("0xEA9D01")
		}
		if !matchError(err, "bad 'count'") {
			t.Error("0xEA33B6")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:0 sn:2 count:1\n"))
		if data != nil {
			t.Error("0xEB21B0")
		}
		if !matchError(err, "out of range") {
			t.Error("0xEB9A96")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:0 sn:1 count:1\n"))
		if data != nil {
			t.Error("0xE11DF3")
		}
		if !matchError(err, "hex") {
			t.Error("0xE43F0E")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:GG sn:1 count:1\n"))
		if data != nil {
			t.Error("0xEF09EC")
		}
		if !matchError(err, "hex") {
			t.Error("0xEF4C9B")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:FF sn:1 count:1\n"))
		if data != nil {
			t.Error("0xE24F86")
		}
		if !matchError(err, "bad hash size") {
			t.Error("0xEB87BE")
		}
	}
	{
		var rc Receiver
		rc.Config = NewDefaultConfig()
		data, err := rc.receiveFragment([]byte(tagFragment +
			"name:abc hash:" +
			"12345678123456781234567812345678" +
			"12345678123456781234567812345678 sn:1 count:1\n"))
		if data != nil {
			t.Error("0xE85E88")
		}
		if !matchError(err, "received no data") {
			t.Error("0xE13A6F")
		}
	}
} //                                              Test_Receiver_receiveFragment_

// (rc *Receiver) sendDataItemHash(req []byte) ([]byte, error)
//
// go test -run Test_Receiver_sendDataItemHash_
//
func Test_Receiver_sendDataItemHash_(t *testing.T) {
	{
		var rc Receiver
		data, err := rc.sendDataItemHash([]byte{})
		if data != nil {
			t.Error("0xE30A2F")
		}
		if !matchError(err, "missing header") {
			t.Error("0xED4B27")
		}
	}
	{
		var rc Receiver
		data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
		if data != nil {
			t.Error("0xED8A3E")
		}
		if !matchError(err, "nil ProvideData") {
			t.Error("0xE65C25")
		}
	}
	{
		var rc Receiver
		rc.ProvideData = func(name string) ([]byte, error) {
			return nil, errors.New("test error")
		}
		data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
		if data != nil {
			t.Error("0xEA5B15")
		}
		if !matchError(err, "test error") {
			t.Error("0xEE3C84")
		}
	}
	{
		var rc Receiver
		rc.ProvideData = func(name string) ([]byte, error) {
			return nil, nil
		}
		data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
		if string(data) != "HASH:"+
			"E3B0C44298FC1C149AFBF4C8996FB924"+
			"27AE41E4649B934CA495991B7852B855" {
			t.Error("0xE8F93C")
		}
		if err != nil {
			t.Error("0xE0D7A2")
		}
	}
	{
		var rc Receiver
		rc.ProvideData = func(name string) ([]byte, error) {
			return []byte("0123456789"), nil
		}
		data, err := rc.sendDataItemHash([]byte(tagDataItemHash))
		if string(data) != "HASH:"+
			"84D89877F0D4041EFB6BF91A16F0248F"+
			"2FD573E6AF05C19F96BEDB9F882F7882" {
			t.Error("0xE37BD7")
		}
		if err != nil {
			t.Error("0xEF13C6")
		}
	}
} //                                             Test_Receiver_sendDataItemHash_

// -----------------------------------------------------------------------------
// # Logging Methods

// (rc *Receiver) logError(id uint32, args ...interface{}) error
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

// (rc *Receiver) logInfo(args ...interface{})
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
