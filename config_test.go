// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[config_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

// to run all tests in this file:
// go test -v -run Test_config_*

// -----------------------------------------------------------------------------

// NewDebugConfig(logFunc ...func(a ...interface{})) *Configuration
//
// go test -run Test_config_NewDebugConfig_
//
func Test_config_NewDebugConfig_(t *testing.T) {
	//
	// returns *Configuration as a string and strips memory addresses
	formatStruct := func(cf *Configuration) string {
		s := fmt.Sprintf("%#v", cf)
		rx := regexp.MustCompile(`\)\(0x.*?\), `)
		ret := string(rx.ReplaceAll([]byte(s), []byte("), ")))
		return ret
	}
	{
		// test!
		got := NewDebugConfig()
		gotS := formatStruct(got)
		//
		// debug configuration should match the one returned
		// by NewDefaultConfig() but with logging activated
		expect := NewDefaultConfig()
		expect.VerboseSender = true
		expect.VerboseReceiver = true
		expect.LogFunc = LogPrint
		expectS := formatStruct(expect)
		//
		if gotS != expectS {
			t.Error("0xEB0A18", "\n",
				"expect:", expectS, "\n",
				"   got:", gotS,
			)
		}
	}
	{
		isFuncCalled := false
		logFunc := func(a ...interface{}) {
			isFuncCalled = true
		}
		// test!
		cf := NewDebugConfig(logFunc)
		cf.LogFunc()
		if !isFuncCalled {
			t.Error("0xEE3AC7")
		}
	}
}

// (cf *Configuration) Validate() error
//
// go test -run Test_config_Configuration_Validate_
//
func Test_config_Configuration_Validate_(t *testing.T) {
	makeValidConfig := func() *Configuration {
		return &Configuration{
			//
			// Components:
			Cipher:     &aesCipher{},
			Compressor: &zlibCompressor{},
			//
			// Limits:
			PacketSizeLimit:   1450,
			PacketPayloadSize: 1024,
			SendBufferSize:    16 * 1024 * 2014, // 16 MiB
			SendRetries:       10,
			//
			// Timeouts and Intervals:
			ReplyTimeout:       15 * time.Second,
			SendPacketInterval: 2 * time.Millisecond,
			SendRetryInterval:  250 * time.Millisecond,
			SendWaitInterval:   50 * time.Millisecond,
			WriteTimeout:       15 * time.Second,
		}
	}
	{
		var cf = makeValidConfig()
		cf.Cipher = nil
		err := cf.Validate()
		if !matchError(err, "nil Configuration.Cipher") {
			t.Error("0xE65F82", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.Compressor = nil
		err := cf.Validate()
		if !matchError(err, "nil Configuration.Compressor") {
			t.Error("0xE2CF8C", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.PacketSizeLimit = 8 - 1
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.PacketSizeLimit") {
			t.Error("0xE6E9BA", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.PacketSizeLimit = (65535 - 8) + 1
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.PacketSizeLimit") {
			t.Error("0xED50EB", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.PacketPayloadSize = 0
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.PacketPayloadSize") {
			t.Error("0xEA00A5", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.PacketSizeLimit = 1000
		cf.PacketPayloadSize = 1000 - 200 + 1
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.PacketPayloadSize") {
			t.Error("0xEC92E8", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.SendBufferSize = -1
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.SendBufferSize") {
			t.Error("0xE2FF75", "wrong error:", err)
		}
	}
	{
		var cf = makeValidConfig()
		cf.SendRetries = -1
		err := cf.Validate()
		if !matchError(err, "invalid Configuration.SendRetries") {
			t.Error("0xE0DE62", "wrong error:", err)
		}
	}
}

// end
