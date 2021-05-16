// -----------------------------------------------------------------------------
// github.com/balacode/udpt                                    /[config_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"fmt"
	"regexp"
	"testing"
)

// to run all tests in this file:
// go test -v -run _config_

// -----------------------------------------------------------------------------

// NewDebugConfig(logFunc ...func(args ...interface{})) *Configuration
//
// go test -run _config_NewDebugConfig_
//
func Test_config_NewDebugConfig_(t *testing.T) {
	//
	// returns *Configuration as a string and strips memory addresses
	formatStruct := func(cfg *Configuration) string {
		s := fmt.Sprintf("%#v", cfg)
		rx := regexp.MustCompile(`\)\(0x.*?\), `)
		ret := string(rx.ReplaceAll([]byte(s), []byte("), ")))
		return ret
	}
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
		t.Error(
			"expect:", expectS, "\n",
			"   got:", gotS,
		)
	}
} //                                                 Test_config_NewDebugConfig_

// end
